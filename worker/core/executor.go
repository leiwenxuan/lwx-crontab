package core

import (
	"context"
	"os/exec"

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
		//jobMan := services.GetJobMangerServer()
		//jobMan.CreateJobLock(info.Job.Name)

		// 执行shell 命令
		cmd = exec.CommandContext(context.TODO(), "C:\\Program Files\\Git\\bin\\bash.exe", "-c", info.Job.Command)
		outPut, err = cmd.CombinedOutput()
		if err != nil {
			logrus.Error("执行任务出错： ", err)
		}
		outPut = outPut
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
