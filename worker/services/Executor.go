package services

var IExecutorServer ExecutorServer

func GetExecutorServer() ExecutorServer {
	return IExecutorServer
}

type ExecutorServer interface {
	InitExecutor() (err error)
}
