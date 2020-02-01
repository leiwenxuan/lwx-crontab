package main

import (
	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
	"github.com/leiwenxuan/crontab/worker"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.EtcdStarter{})
	infra.Register(&worker.EtcdRegisterStarter{})
	infra.Register(&base.MongoDBStarter{})
	//infra.Register(&worker.WatchRegisterStarter{})
	infra.Register(&base.HookStarter{})

}
