package core

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"go.etcd.io/etcd/clientv3"
)

type JobLock struct {
	Kv    clientv3.KV
	Lease clientv3.Lease

	JobName    string             // 任务名字
	cancelFunc context.CancelFunc // 用于终止自动续租
	leaseId    clientv3.LeaseID   // 租约ID
	isLocked   bool               // 是否上锁成功
}

func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		Kv:      kv,
		Lease:   lease,
		JobName: jobName,
	}
	return
}

func (j *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)

	// 1, 创建租约(5秒)
	if leaseGrantResp, err = j.Lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println("Lease.Grant: ", err)

		return
	}

	// context用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	// 续租ID
	leaseId = leaseGrantResp.ID

	// 2, 自动续租
	if keepRespChan, err = j.Lease.KeepAlive(cancelCtx, leaseId); err != nil {
		fmt.Println("KeepAlive: ", err)
		goto FAIL
	}

	// 3, 处理续租应答的协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan: // 自动续租的应答
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()
	// 创建事务
	txn = j.Kv.Txn(context.TODO())
	// 锁路径
	lockKey = JOB_LOCK_DIR + j.JobName

	// 5,事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
		fmt.Println("txn.Commit(): ", err)

	}

	// 6,成功返回，失败释放租约
	if !txnResp.Succeeded {
		// 锁被占用
		err = ERR_LOCK_ALREADY_REQUIRED
		fmt.Println("!txnResp.Succeeded ", err)

		goto FAIL
	}
	// 抢锁成功
	j.leaseId = leaseId
	j.cancelFunc = cancelFunc
	j.isLocked = true
	return
FAIL:
	cancelFunc()
	logrus.Error("取消自动续租协程")
	// 取消自动续租协程
	_, _ = j.Lease.Revoke(context.TODO(), j.leaseId) // 释放租约
	return err
}

// 释放锁
func (j *JobLock) Unlock() {
	if j.isLocked {
		j.cancelFunc()
		_, _ = j.Lease.Revoke(context.TODO(), j.leaseId)
	}
}
