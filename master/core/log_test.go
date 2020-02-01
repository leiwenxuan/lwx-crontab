package core

import (
	"fmt"
	"testing"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
	"github.com/tietang/props/ini"
)

func TestLogServer_ListLog(t *testing.T) {

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.MongoDBStarter{})

	//file := kvs.GetCurrentFilePath("F:\\code\\002Golang\\lwx-crontab\\master\\main\\config.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource("F:\\code\\002Golang\\lwx-crontab\\master\\main\\config.ini")
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
	log := new(LogServer)
	result, count, err := log.ListLog("work_03", 0, 60)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(len(result), count)
}
