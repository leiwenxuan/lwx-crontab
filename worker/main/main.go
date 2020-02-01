package main

import (
	"runtime"
	"time"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
	_ "github.com/leiwenxuan/crontab/worker"
	_ "github.com/leiwenxuan/crontab/worker/core"
	"github.com/leiwenxuan/crontab/worker/services"
	"github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

// 初始化线程数量
func initEnv() {
	logrus.Info("初始化线程数量：", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// 初始化线程
	initEnv()
	//获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("config.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
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
