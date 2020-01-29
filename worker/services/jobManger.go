package services

var IJobMangerServer JobMangerServer

func GetJobMangerServer() JobMangerServer {
	return IJobMangerServer
}

type JobMangerServer interface {
	// 监听任务变化
	JobWatch() (err error)
	// 监听强杀任务通知
	WatchKiller()
	// 创建任务执行锁
	CreateJobLock(jobName string) (jobLock *JobLock)
}
