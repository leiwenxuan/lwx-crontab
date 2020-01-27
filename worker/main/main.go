package main

import (
	"runtime"
	"time"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
	_ "github.com/leiwenxuan/crontab/worker"
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
	time.Sleep(2 * time.Second)
	//服务注册，注册本机IP到etcd，并不断续租
	//err := worker.InitRegister()
	//logrus.Println(err)
}
