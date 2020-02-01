package core

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/leiwenxuan/crontab/worker/services"
)

var _ services.SchedulerServer = new(SchedulerJobSer)

type SchedulerSer struct {
	jobEventChan      chan *services.JobEvent              // etcd任务事件队列
	jobPlanTable      map[string]*services.JobSchedulePlan // 任务调度计划表
	jobExecutingTable map[string]*services.JobExecuteInfo  // 任务执行表
	jobResultChan     chan *services.JobExecuteResult      // 任务结果队列
}

var (
	Gscheduler *SchedulerSer
)

type SchedulerJobSer struct {
}

func (s SchedulerJobSer) InitScheduler() (err error) {
	Gscheduler = &SchedulerSer{
		jobEventChan:      make(chan *services.JobEvent, 1000),
		jobPlanTable:      make(map[string]*services.JobSchedulePlan),
		jobExecutingTable: make(map[string]*services.JobExecuteInfo),
		jobResultChan:     make(chan *services.JobExecuteResult, 1000),
	}

	// 启动调度协程
	go Gscheduler.ScheduleLoop()
	return
}

var onceScheduler sync.Once

func init() {
	onceScheduler.Do(func() {
		services.ISchedulerServer = new(SchedulerJobSer)
		logrus.Debug("调度器初始化")
		_ = services.GetSchedulerServer().InitScheduler()
	})
}

func (s *SchedulerSer) HandleJobEvent(jobEvent *services.JobEvent) {
	var (
		jobSchedulePlan *services.JobSchedulePlan
		jobExecuteInfo  *services.JobExecuteInfo
		jobExisted      bool
		err             error
		jobExecuting    bool
	)
	switch jobEvent.EventType {
	case JOB_EVENT_SAVE:
		// 保存事件
		if jobSchedulePlan, err = services.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		s.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case JOB_EVENT_DELETE:
		// 删除事件
		if jobSchedulePlan, jobExisted = s.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(s.jobPlanTable, jobEvent.Job.Name)
		}
	case JOB_EVENT_KILL:
		if jobExecuteInfo, jobExecuting = s.jobExecutingTable[jobEvent.Job.Name]; jobExecuting {
			// 触发命令执行， 杀死进程
			jobExecuteInfo.CancelFunc()
		}
	}

}

func (s *SchedulerSer) TryStartJob(jobPlan *services.JobSchedulePlan) {
	//调度和执行是2件事情
	var (
		jobExecuterInfo *services.JobExecuteInfo
		jobExecuting    bool
	)
	// 执行的任务可能运行很久，1分钟会调度60次，但是只能执行1次，防止并发!

	// 如果任务正在执行，跳过本次调度
	if jobExecuterInfo, jobExecuting = s.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		logrus.Infof("任务尚未推出: %", jobPlan.Job.Name)
		return
	}

	// 构建执行状态
	jobExecuterInfo = services.BuildJobExecuteInfo(jobPlan)

	// 保存执行状态
	s.jobExecutingTable[jobPlan.Job.Name] = jobExecuterInfo

	// 执行任务
	logrus.Info("执行任务：  ", jobExecuterInfo.Job.Name, jobExecuterInfo.PlanTime, jobExecuterInfo.RealTime)
	// TODO 执行任务
	G_executor.ExecutorJob(jobExecuterInfo)

}

func (s *SchedulerSer) TrySchedule() (schdulerAfter time.Duration) {
	var (
		jobPlan  *services.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)
	// 如果任务为空的， 睡眠1秒
	if len(s.jobPlanTable) == 0 {
		schdulerAfter = 1 * time.Second
		return
	}
	// 当前事件
	now = time.Now()
	// 遍历所有任务
	for _, jobPlan = range s.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			logrus.Info("执行任务: ", jobPlan.Job.Name)
			s.TryStartJob(jobPlan)
			// 更新下次执行任务
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}
		// 统计最近一个要过期的任务
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度间隔(最近时间-当前时间)
	schdulerAfter = (*nearTime).Sub(now)
	return
}

func (s *SchedulerSer) HandleJobResult(result *services.JobExecuteResult) {

	var jobLog *services.JobLog
	delete(s.jobExecutingTable, result.ExecuteInfo.Job.Name)
	//if result.Err  != ERR_LOCK_ALREADY_REQUIRED
	// 生产日志
	if result.Err != ERR_LOCK_ALREADY_REQUIRED {
		jobLog = &services.JobLog{
			JobName:      result.ExecuteInfo.Job.Name,
			Command:      result.ExecuteInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
		GLogSink.Append(jobLog)
	}
	logrus.Info("任务执行完成:  ", result.ExecuteInfo.Job.Name)

	return
}

func (s *SchedulerSer) PushJobEvent(jobEvent *services.JobEvent) {
	logrus.Debug("PushJobEvent", jobEvent)

	s.jobEventChan <- jobEvent
}

func (s *SchedulerSer) PushJobResult(jobResult *services.JobExecuteResult) {
	logrus.Debug("PushJobResult")

	s.jobResultChan <- jobResult
}

func (s *SchedulerSer) ScheduleLoop() {
	var (
		jobEvent      *services.JobEvent
		scheduleAfter time.Duration
		scheduleTImer *time.Timer
		jobResult     *services.JobExecuteResult
	)
	// 初始化一次
	scheduleAfter = s.TrySchedule()

	// 调度的延迟时间
	scheduleTImer = time.NewTimer(scheduleAfter)

	// 定时任务
	for {
		select {
		case jobEvent = <-s.jobEventChan:
			// 维护的任务队列
			s.HandleJobEvent(jobEvent)
		case <-scheduleTImer.C:
		case jobResult = <-s.jobResultChan:
			// 监听任务执行结果
			s.HandleJobResult(jobResult)
		}
		// 调度任务
		scheduleAfter = s.TrySchedule()
		scheduleTImer.Reset(scheduleAfter)
	}
}
