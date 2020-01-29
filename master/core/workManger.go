package core

import (
	"context"
	"sync"

	"github.com/leiwenxuan/crontab/infra/base"
	"github.com/leiwenxuan/crontab/master/services"
	"go.etcd.io/etcd/clientv3"
)

var _ services.WorkerMangerServer = new(WorkServer)

type WorkServer struct {
}

var onceWork sync.Once

func init() {
	onceWork.Do(func() {
		services.IWorkerServer = new(WorkServer)
	})
}

func (w WorkServer) WorkerList() (workerArr []string, err error) {
	var gerResp *clientv3.GetResponse
	client := base.EtcdClient()
	kv := clientv3.NewKV(client)

	if gerResp, err = kv.Get(context.TODO(), JOB_WORKER_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	// 解析每个节点
	for _, kv := range gerResp.Kvs {
		workerIp := ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, workerIp)
	}
	return
}
