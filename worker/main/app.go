package main

import (
	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.EtcdStarter{})
	//infra.Register(&services.EtcdRegisterStarter{})
	infra.Register(&base.MongoDBStarter{})
	//infra.Register(&worker.WatchRegisterStarter{})
	infra.Register(&base.HookStarter{})

}
