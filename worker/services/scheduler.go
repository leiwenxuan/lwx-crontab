package services

var ISchedulerServer SchedulerServer

func GetSchedulerServer() SchedulerServer {
	return ISchedulerServer
}

type SchedulerServer interface {
	// 处理任务事件
	InitScheduler() (err error)
}
