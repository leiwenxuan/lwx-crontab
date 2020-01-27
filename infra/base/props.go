package base

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"github.com/leiwenxuan/crontab/infra"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	Check(props)
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	log.Info("初始化配置.")
}
