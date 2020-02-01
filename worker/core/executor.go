package core

import (
	"context"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/leiwenxuan/crontab/worker/services"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecutorJob(info *services.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			outPut  []byte
			result  *services.JobExecuteResult
			jobLock *JobLock
		)

		result = &services.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		// 初始化分布式锁
		logrus.Debug("初始化分布式锁KEY: ", info.Job.Name)
		jobLock = GJobServer.CreateJobLock(info.Job.Name)
		// 记录任务开始时间
		result.StartTime = time.Now()

		// 上锁
		//time.Sleep(1 * time.Millisecond)
		defer jobLock.Unlock()

		err = jobLock.TryLock()

		if err != nil {
			// 上锁失败
			result.Err = err
			result.EndTime = time.Now()
			logrus.Debug("加锁失败")

		} else {
			logrus.Debug("加锁成功")
			result.StartTime = time.Now()
			// 执行shell 命令
			//cmd = exec.CommandContext(context.TODO(), "C:\\Program Files\\Git\\bin\\bash.exe", "-c", info.Job.Command)
			cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
			outPut, err = cmd.CombinedOutput()
			if err != nil {
				logrus.Error("执行任务出错： ", err)
			}
			result.EndTime = time.Now()
			result.Output = outPut
			result.Err = err
		}

		// 任务执行完成后，把执行的结果返回给Scheduler,Scheduler会从executingTable中删除掉执行记录
		Gscheduler.HandleJobResult(result)

	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
func init() {
	_ = InitExecutor()
}
