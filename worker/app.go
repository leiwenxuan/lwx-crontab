package worker

import (
	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.EtcdStarter{})
	infra.Register(&base.EtcdRegisterStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.HookStarter{})
}
