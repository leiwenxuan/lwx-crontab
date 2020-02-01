package core

import (
	"errors"
	"strings"
)

const (
	// 任务保存目录
	JOB_SAVE_DIR = "/cron/jobs/"

	// 任务强杀目录
	JOB_KILLER_DIR = "/cron/killer/"

	// 任务锁目录
	JOB_LOCK_DIR = "/cron/lock/"

	// 服务注册目录
	JOB_WORKER_DIR = "/cron/workers/"

	// 保存任务事件
	JOB_EVENT_SAVE = 1

	// 删除任务事件
	JOB_EVENT_DELETE = 2

	// 强杀任务事件
	JOB_EVENT_KILL      = 3
	JobLogCommitTimeout = 5000
	JobLogBatchSize     = 100
)

var (
	ERR_LOCK_ALREADY_REQUIRED = errors.New("锁已经被占用")
	ERR_NO_LOCAL_IP_FOUND     = errors.New("没有找到网卡IP")
)

// 提取worker的IP
func ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}

// 从/cron/killer/job10提取job10
func ExtractKillerName(killerKey string) string {
	return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

// 从etcd的key中提取任务名
// /cron/jobs/job10抹掉/cron/jobs/
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}
