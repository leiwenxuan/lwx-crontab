package apis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
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
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	groupRouter := base.Iris().Party("/v1/master", crs).AllowMethods(iris.MethodOptions)
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

	// 登录用户
	groupRouter.Post("/login", m.LoginUser)
	// 菜单栏

	groupRouter.Get("/menus", m.Menus)
}

type LoginUser struct {
	Id       int    `json:"id"`
	Rid      int    `json:"rid"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type PurviewMenus struct {
	ID       int         `json:"id"`
	AuthName string      `json:"authName"`
	Path     interface{} `json:"path"`
	Children []Children  `json:"children"`
}
type Children struct {
	ID       int           `json:"id"`
	AuthName string        `json:"authName"`
	Path     interface{}   `json:"path"`
	Children []interface{} `json:"children"`
}

func (m *MasterApi) Menus(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	var menusList = ` [
        {
            "id": 125,
            "authName": "健康检查",
            "path": "worker",
            "children": [
                {
                    "id": 110,
                    "authName": "Worker节点",
                    "path": "worker",
                    "children": [],
                    "order": null
                }
            ],
            "order": 1
        },
        {
            "id": 103,
            "authName": "定时任务",
            "path": "job",
            "children": [
                {
                    "id": 111,
                    "authName": "任务列表",
                    "path": "job",
                    "children": [],
                    "order": null
                }
            ],
            "order": 2
        }
    ]`

	var menusWork []*PurviewMenus
	if err := json.Unmarshal([]byte(menusList), &menusWork); err != nil {
		logrus.Error("转码失败", err)
	}
	//var oneMenusWork = &PurviewMenus{
	//	ID:       100,
	//	AuthName: "健康节点",
	//	Path:     "/worker/list",
	//	Children: nil,
	//}
	//r.Data = []interface{}{oneMenusWork}
	r.Data = menusWork

	r.Code = 200
	_, _ = ctx.JSON(&r)
	logrus.Infof("[/login]成功:%s\t%+v\n", time.Now().Format("2006-01-02 15:04:05"))

}

func (m *MasterApi) LoginUser(ctx iris.Context) {
	r := base.Res{
		Code: base.ResCodeOk,
	}
	var login = &LoginUser{
		Id:       100,
		Rid:      0,
		Username: "雷文轩",
		Mobile:   "00000000",
		Email:    "00000",
		Token:    "token",
	}
	r.Data = login
	r.Code = 200
	_, _ = ctx.JSON(&r)
	logrus.Infof("[/login]成功:%s\t%+v\n", time.Now().Format("2006-01-02 15:04:05"))

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
	var workListDict []interface{}
	for _, v := range result {
		var workDict = map[string]string{
			"workerIp": "",
		}
		workDict["workerIp"] = v
		workListDict = append(workListDict, workDict)
	}
	r.Data = workListDict
	_, _ = ctx.JSON(&r)
	logrus.Infof("[/worker/list]成功:%s\t%d\n", time.Now().Format("2006-01-02 15:04:05"), len(result))

}

// 日志列表
func (m *MasterApi) LogList(ctx iris.Context) {
	var err error
	logParam := services.LogParam{}
	logParam.Name = ctx.URLParamDefault("name", "test")
	logParam.Limit = ctx.URLParamInt64Default("limit", 30)

	logParam.Skip = ctx.URLParamInt64Default("skip", 0)
	r := base.Res{Code: base.ResCodeOk}
	fmt.Println("m.logSer", m.logSer)
	result, count, err := m.logSer.ListLog(logParam.Name, logParam.Skip, logParam.Limit)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		logrus.Error(err)
	}
	var data = map[string]interface{}{
		"count": count,
		"data":  result,
	}

	r.Data = data
	_, _ = ctx.JSON(&r)

}
