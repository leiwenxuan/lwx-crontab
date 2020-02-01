package services

type LogParam struct {
	Name  string `json:"name"`
	Limit int64  `json:"limit"`
	Skip  int64  `json:"skip"`
}

var ILoggerMangerServer LoggerMongerServer

func GetLoggerMangerServer() LoggerMongerServer {

	return ILoggerMangerServer
}

type LoggerMongerServer interface {
	// 查看日志
	ListLog(taskName string, skip int64, limit int64) (jobLog []JobLog, count int64, err error)
}
