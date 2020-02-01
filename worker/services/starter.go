package services

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/leiwenxuan/crontab/infra"
)

type WatchRegisterStarter struct {
	infra.BaseStarter
}

func (s *WatchRegisterStarter) Start(ctx infra.StarterContext) {
	jobManger := GetJobMangerServer()
	// 初始化job管理器
	_ = jobManger.InitJobManger()
	// 初始化日志
	logManger := GetLogManger()
	_ = logManger.InitLogSink()
	// 服务注册
	if err := InitRegister(); err != nil {
		logrus.Error("服务注册失败", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
