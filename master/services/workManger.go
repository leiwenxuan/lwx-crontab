package services

var IWorkerServer WorkerMangerServer

func GetWorkerServer() WorkerMangerServer {
	return IWorkerServer
}

type WorkerMangerServer interface {
	WorkerList() (workerArr []string, err error)
}
