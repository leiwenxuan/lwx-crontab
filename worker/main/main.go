package main

import (
	"runtime"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"

	//_ "github.com/leiwenxuan/crontab/worker"
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
	logrus.Debug("调度器初始化")
	_ = services.GetSchedulerServer().InitScheduler()
	// 初始化执行器
	logrus.Debug("初始化执行器")
	_ = services.GetExecutorServer().InitExecutor()

	// 初始化job管理器
	logrus.Debug("初始化job管理器")
	_ = services.GetJobMangerServer().InitJobManger()
	_ = services.InitRegister()
	// 初始化日志
	_ = services.GetLogManger().InitLogSink()
	logrus.Debug("初始化日志")

	ch := make(chan int, 1)
	<-ch

}
