package test

import (
	"testing"

	"github.com/leiwenxuan/crontab/worker"
)

func TestGetIp(t *testing.T) {
	worker.GetLocalIP()
}
