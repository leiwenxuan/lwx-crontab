package core

import (
	"fmt"
	"testing"

	"github.com/leiwenxuan/crontab/master/services"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
	"github.com/tietang/props/ini"
)

func TestJobServer_SaveJob(t *testing.T) {

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.EtcdStarter{})

	//file := kvs.GetCurrentFilePath("F:\\code\\002Golang\\lwx-crontab\\master\\main\\config.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource("F:\\code\\002Golang\\lwx-crontab\\master\\main\\config.ini")
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
	log := new(JobServer)

	var job = services.Job{
		Name:     "work_2",
		Command:  "test",
		CronExpr: "* * * * * echo 'hell word'",
	}
	resultsave, err := log.SaveJob(&job)
	if err != nil {
		fmt.Println("resultsave", err)
	}
	fmt.Println(resultsave)

	resultdel, err := log.DeleteJob("work_1")
	if err != nil {
		fmt.Println("resultdel", resultdel)
	}
	fmt.Println("resultdel", resultdel)
	getresult, err := log.ListJob()
	if err != nil {
		fmt.Println("err getresult", err)
	}

	fmt.Println("get list ")

}
