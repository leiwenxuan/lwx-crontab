package apis

import (
	"fmt"
	"time"

	"github.com/kataras/iris"
	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
	"github.com/leiwenxuan/crontab/master/services"
	"github.com/sirupsen/logrus"
)

func init() {
	infra.RegisterApi(new(MasterApi))

}

type MasterApi struct {
	logSer  services.LoggerMongerServer
	jobSer  services.JobMangerServer
	workSer services.WorkerMangerServer
}

func (m *MasterApi) Init() {
	m.logSer = services.GetLoggerMangerServer()
	m.jobSer = services.GetJobMangerServer()
	m.workSer = services.GetWorkerServer()
	groupRouter := base.Iris().Party("/v1/master")
	// 日志
	groupRouter.Get("/job/log", m.LogList)
	// 任务
	// 1.1 保存任务
	groupRouter.Post("/job/save", m.JobSave)
	// 1.2 删除任务
	groupRouter.Post("/job/delete", m.JobDelete)
	// 1.3 列出任务
	groupRouter.Get("/job/list", m.JobList)
	// 1.4  杀掉任务
	groupRouter.Post("/job/kill", m.JobKill)

	// 3.1 健康节点
	groupRouter.Get("/worker/list", m.WorkList)
}

type DeleteJob struct {
	Name string `json:"name"`
}

//  保存任务
func (m *MasterApi) JobSave(ctx iris.Context) {
	job := services.Job{}
	err := ctx.ReadJSON(&job)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	result, err := m.jobSer.SaveJob(&job)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r.Data = result
	r.Message = "success"
	_, _ = ctx.JSON(&r)
	logrus.Infof("[/job/save]成功:%s\t%+v\n", time.Now().Format("2006-01-02 15:04:05"), job)
}

// 删除任务
func (m *MasterApi) JobDelete(ctx iris.Context) {
	var d = DeleteJob{}
	err := ctx.ReadJSON(&d)
	r := base.Res{
		Code:    base.ResCodeOk,
		Message: "",
		Data:    nil,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	result, err := m.jobSer.DeleteJob(d.Name)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	r.Data = result
	_, _ = ctx.JSON(&r)
	logrus.Infof("[/job/delete]成功:%s\t%+v\n", time.Now().Format("2006-01-02 15:04:05"), result)

}

// 列出任务
func (m *MasterApi) JobList(ctx iris.Context) {
	r := base.Res{
		Code:    base.ResCodeOk,
		Message: "",
		Data:    nil,
	}
	result, err := m.jobSer.ListJob()
	if err != nil {
		r.Message = err.Error()
		r.Code = base.ResCodeInnerServerError
		ctx.JSON(&r)
		logrus.Error("JobList", err)
	}
	r.Data = result
	_, _ = ctx.JSON(&r)
	logrus.Infof("[/job/list]成功:%s\t%+v\n", time.Now().Format("2006-01-02 15:04:05"))

}

type KillJob struct {
	Name string `json:"name"`
}

// 杀掉任务
func (m *MasterApi) JobKill(ctx iris.Context) {
	r := base.Res{
		Code:    0,
		Message: "",
		Data:    nil,
	}
	var k = KillJob{}
	err := ctx.ReadJSON(&k)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error("err ReadJson: ", err)
	}
	err = m.jobSer.KillJob(k.Name)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error("err KillJob: ", err)
	}
	_, _ = ctx.JSON(&r)
}

// worker 节点
func (m *MasterApi) WorkList(ctx iris.Context) {
	r := base.Res{
		Code:    base.ResCodeOk,
		Message: "",
		Data:    nil,
	}
	result, err := m.workSer.WorkerList()
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error("err m.workSer.WorkerList(): ", err)
	}

	r.Data = result
	_, _ = ctx.JSON(&r)
	logrus.Infof("[/worker/list]成功:%s\t%d\n", time.Now().Format("2006-01-02 15:04:05"), len(result))

}

// 日志列表
func (m *MasterApi) LogList(ctx iris.Context) {
	var err error
	logParam := services.LogParam{}
	logParam.Name = ctx.URLParamDefault("name", "test")
	logParam.Limit = ctx.URLParamInt64Default("limit", 1)
	if err != nil {
		logrus.Error("ctx.URLParamInt64(limit)", err)
	}
	logParam.Skip = ctx.URLParamInt64Default("skip", 0)
	if err != nil {
		logrus.Error("ctx.URLParamInt64(skip)", err)
	}
	r := base.Res{Code: base.ResCodeOk}
	fmt.Println("m.logSer", m.logSer)
	result, err := m.logSer.ListLog(logParam.Name, logParam.Skip, logParam.Skip)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		logrus.Error(err)
	}
	r.Data = result
	_, _ = ctx.JSON(&r)

}
