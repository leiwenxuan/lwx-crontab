package core

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/leiwenxuan/crontab/infra/base"
	"github.com/leiwenxuan/crontab/worker/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogSink struct {
	Client         *mongo.Client
	LogCollection  *mongo.Collection
	LogChan        chan *services.JobLog
	AutoCommitChan chan *services.LogBatch
}

var _ services.LogMangerServer = new(LogServer)

type LogServer struct {
}

var logOnce sync.Once

func init() {
	logOnce.Do(func() {
		services.ILogMangerServer = new(LogServer)

	})
}

var (
	// 单利
	GLogSink *LogSink
)

func (l LogServer) InitLogSink() (err error) {
	client := base.ClientMongodb()
	GLogSink = &LogSink{
		Client:         client,
		LogCollection:  client.Database("cron").Collection("log"),
		LogChan:        make(chan *services.JobLog, 1000),
		AutoCommitChan: make(chan *services.LogBatch, 1000),
	}
	go GLogSink.WriteLoop()
	return err
}

// 程序配置
type Config struct {
	EtcdEndpoints         []string `json:"etcdEndpoints"`         // etcd集群地址
	EtcdDialTimeout       int      `json:"etcdDialTimeout"`       // etcd超时时间
	MongodbUri            string   `json:"mongodbUri"`            // mongo地址
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"` // mongo连接超时时间
	JobLogBatchSize       int      `json:"jobLogBatchSize"`       // 日志批量大小
	JobLogCommitTimeout   int      `json:"jobLogCommitTimeout"`   // 日志自动提交超时时间
}

// 日志存储协程
func (l *LogSink) WriteLoop() {
	var (
		log          *services.JobLog
		logBatch     *services.LogBatch
		commitTimer  *time.Timer
		timeoutBatch *services.LogBatch
	)
	for {
		select {
		case log = <-l.LogChan:
			if logBatch == nil {
				logBatch = &services.LogBatch{}
				// 让这个批次超时自动提交(给1秒的时间)
				commitTimer = time.AfterFunc(
					time.Duration(JobLogCommitTimeout)*time.Millisecond,
					func(batch *services.LogBatch) func() {
						return func() {
							// 超时提交
							l.AutoCommitChan <- batch
						}
					}(logBatch),
				)
			}
			// 把更新的日志写入批次中
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果批次满了
			if len(logBatch.Logs) >= JobLogBatchSize {
				// 发送日志
				l.SaveLogs(logBatch)
				// 清空logbatch
				logBatch = nil
				// 取消定时
				commitTimer.Stop()
			}
		case timeoutBatch = <-l.AutoCommitChan:
			// 过期的批次
			// 判断过期的， 是否是当前的
			if timeoutBatch != logBatch {
				continue
			}
			// 把批次写入mongo
			l.SaveLogs(timeoutBatch)
			logBatch = nil
		}
	}

}

// 批量写入
func (l *LogSink) SaveLogs(batch *services.LogBatch) {
	l.LogCollection.InsertMany(context.TODO(), batch.Logs)
	logrus.Debug("提交日志:  ", time.Now().Format("2006-01-02 15:04:05"))
}

// 发送日志
func (l *LogSink) Append(jobLog *services.JobLog) {
	logrus.Info("发送日志")
	select {
	case l.LogChan <- jobLog:
	default:
		// 队里满了就丢弃

	}

}
