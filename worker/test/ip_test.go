package test

import (
	"testing"

	"github.com/leiwenxuan/crontab/worker/services"
)

func TestGetIp(t *testing.T) {
	services.GetLocalIP()
}
