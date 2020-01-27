package services

var IJobMangerServer JobMangerServer

func GetJobMangerServer() JobMangerServer {
	return IJobMangerServer
}

type JobMangerServer interface {
	// 保存任务
	SaveJob(job *Job) (oldJob *Job, err error)
	// 删除任务
	DeleteJob(name string) (oldJob *Job, err error)
	// 列举任务
	ListJob() (jobList []*Job, err error)
	// 杀死任务
	KillJob(name string) (err error)
}
