package services

var IJobMangerServer JobMangerServer

func GetJobMangerServer() JobMangerServer {
	return IJobMangerServer
}

type JobMangerServer interface {
	// 初始化
	InitJobManger() (err error)
}
