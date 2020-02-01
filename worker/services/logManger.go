package services

var ILogMangerServer LogMangerServer

func GetLogManger() LogMangerServer {
	return ILogMangerServer
}

type LogMangerServer interface {
	// 初始化日志模块
	InitLogSink() (err error)
}
