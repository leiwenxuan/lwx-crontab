package core

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/leiwenxuan/crontab/infra/base"
	"github.com/leiwenxuan/crontab/worker/services"
	"go.etcd.io/etcd/clientv3"
	//"go.etcd.io/etcd/mvcc/mvccpb"
)

var _ services.JobMangerServer = new(JobServer)

type JobMangerSer struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	Watcher clientv3.Watcher
}

type JobServer struct {
}

var (
	// 单利 模式
	GJobServer *JobMangerSer
)

func (j JobServer) InitJobManger() (err error) {
	client := base.EtcdClient()
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	watcher := clientv3.NewWatcher(client)

	GJobServer = &JobMangerSer{
		Client:  client,
		Kv:      kv,
		Lease:   lease,
		Watcher: watcher,
	}

	_ = GJobServer.JobWatch()
	GJobServer.WatchKiller()
	return
}

var onceJob sync.Once

func init() {
	onceJob.Do(func() {
		services.IJobMangerServer = new(JobServer)
	})
}

func (j *JobMangerSer) JobWatch() (err error) {
	var (
		getResp  *clientv3.GetResponse
		job      *services.Job
		jobEvent *services.JobEvent
	)

	//var schedulerJob = services.GetSchedulerServer()

	client := base.EtcdClient()
	kv := clientv3.NewKV(client)
	jobWatch := clientv3.NewWatcher(client)
	// 获取所有的job任务
	if getResp, err = kv.Get(context.TODO(), JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	logrus.Info("all kv:   ", getResp)
	// 当前任务
	for _, kvpair := range getResp.Kvs {
		if job, err = services.UnpackJob(kvpair.Value); err == nil {
			jobEvent := services.BuildJobEvent(JOB_EVENT_SAVE, job)
			// 同步给调度协程
			logrus.Debug(jobEvent)
			Gscheduler.PushJobEvent(jobEvent)
		}
	}
	// 从当前revision 版本向后监听
	go func() {
		watchStartRevision := getResp.Header.Revision
		// 监听job 目录变化
		watchChan := jobWatch.Watch(context.TODO(), JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

		// 监听处理事件
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					// 1 保存shijian
					if job, err = services.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 构建event
					jobEvent = services.BuildJobEvent(JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					jobName := ExtractJobName(string(watchEvent.Kv.Key))

					job = &services.Job{Name: jobName}
					jobEvent = services.BuildJobEvent(JOB_EVENT_DELETE, job)
				}
				logrus.Infof("监听事件变化: \t %d , %+v", jobEvent.EventType, jobEvent.Job)
				// TODO 推送给调度协程
				Gscheduler.PushJobEvent(jobEvent)
			}
		}

	}()
	return
}

func (j *JobMangerSer) WatchKiller() {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent   *services.JobEvent
		jobName    string
		job        *services.Job
	)
	client := base.EtcdClient()
	jobWatch := clientv3.NewWatcher(client)
	// 监听killer
	go func() {
		watchChan = jobWatch.Watch(context.TODO(), JOB_KILLER_DIR, clientv3.WithPrefix())
		// 监听处理事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					jobName = ExtractKillerName(string(watchEvent.Kv.Key))
					job = &services.Job{
						Name: jobName,
					}
					jobEvent = services.BuildJobEvent(JOB_EVENT_KILL, job)
					Gscheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE:
					// killer标记过期，被自动删除

				}
			}
		}
	}()

}

func (j *JobMangerSer) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, j.Kv, j.Lease)
	return jobLock
}
