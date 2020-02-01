package core

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/leiwenxuan/crontab/infra/base"
	"go.etcd.io/etcd/clientv3"

	"github.com/leiwenxuan/crontab/master/services"
)

var _ services.JobMangerServer = new(JobServer)

type JobServer struct {
}

var onceJob sync.Once

func init() {
	onceJob.Do(func() {
		services.IJobMangerServer = new(JobServer)
	})
}
func (j JobServer) SaveJob(job *services.Job) (oldJob *services.Job, err error) {
	var (
		client    *clientv3.Client
		kv        clientv3.KV
		putResp   *clientv3.PutResponse
		jobJson   []byte
		oldJobObj services.Job
	)
	// 任务名称
	jobKey := JOB_SAVE_DIR + job.Name
	// 任务信息json格式
	if jobJson, err = json.Marshal(job); err != nil {
		return
	}
	fmt.Println(string(jobJson))
	client = base.EtcdClient()
	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), jobKey, string(jobJson), clientv3.WithPrevKV()); err != nil {
		return
	}
	// 判断是更新还是新增
	if putResp.PrevKv != nil {
		// 对旧值做反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

func (j JobServer) DeleteJob(name string) (oldJob *services.Job, err error) {
	var (
		delResp   *clientv3.DeleteResponse
		oldJobObj services.Job
	)
	jobKey := JOB_SAVE_DIR + name
	fmt.Println("jobKey: ", jobKey)
	client := base.EtcdClient()
	kv := clientv3.NewKV(client)
	if delResp, err = kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	fmt.Println("len:  ", len(delResp.PrevKvs))
	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}

	return oldJob, err
}

func (j JobServer) ListJob() (jobList []*services.Job, err error) {
	var (
		getResp *clientv3.GetResponse
		dirkey  string
	)
	dirkey = JOB_SAVE_DIR

	client := base.EtcdClient()
	kv := clientv3.NewKV(client)
	if getResp, err = kv.Get(context.TODO(), dirkey, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, kvPair := range getResp.Kvs {
		job := &services.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

func (j JobServer) KillJob(name string) (err error) {
	var (
		killKey        string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
	)

	killKey = JOB_KILLER_DIR + name
	logrus.Info("killKey： ", killKey)
	client := base.EtcdClient()
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	if leaseGrantResp, err = lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	// 租约
	leaseId = leaseGrantResp.ID
	logrus.Info("leaseId: ", leaseId)
	// 设置killer 标记
	if _, err = kv.Put(context.TODO(), killKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}
	return
}
