package main

import (
	"time"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/worker/services"
)

type WatchRegisterStarter struct {
	infra.BaseStarter
}

func (s *WatchRegisterStarter) Start(ctx infra.StarterContext) {
	jobManger := services.GetJobMangerServer()
	// 初始化job管理器
	_ = jobManger.InitJobManger()
	// 初始化日志
	logManger := services.GetLogManger()
	_ = logManger.InitLogSink()
	for {
		time.Sleep(1 * time.Second)
	}
}
