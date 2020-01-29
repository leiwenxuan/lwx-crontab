package worker

import (
	"time"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/leiwenxuan/crontab/worker/services"
)

type WatchRegisterStarter struct {
	infra.BaseStarter
}

func (s *WatchRegisterStarter) Start(ctx infra.StarterContext) {
	jobSer := services.GetJobMangerServer()
	_ = jobSer.JobWatch()
	jobSer.WatchKiller()
	for {
		time.Sleep(1 * time.Second)
	}
}
