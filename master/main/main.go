package main

import (
	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
	_ "github.com/leiwenxuan/crontab/master"
	_ "github.com/leiwenxuan/crontab/master/apis"
	_ "github.com/leiwenxuan/crontab/master/core"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	//获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("config.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
}
