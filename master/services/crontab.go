package services

import (
	"context"
	"time"

	"github.com/gorhill/cronexpr"
)

type Job struct {
	Name     string `json:"name"`     // 任务名
	Command  string `json:"command"`  // shell命令
	CronExpr string `json:"cronExpr"` // cron表达式
}

// 任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 // 要调度的任务信息
	Expr     *cronexpr.Expression // 解析好的cronexpr表达式
	NextTime time.Time            // 下次调度时间
}

// 任务执行状态
type JobExecuteInfo struct {
	Job        *Job               // 任务信息
	PlanTime   time.Time          // 理论上的调度时间
	RealTime   time.Time          // 实际的调度时间
	CancelCtx  context.Context    // 任务command的context
	CancelFunc context.CancelFunc // 用于取消command执行的cancel函数
}

// HTTP接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 变化事件
type JobEvent struct {
	EventType int // Save, Delete
	Job       *Job
}

// 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          // 脚本输出
	Err         error           // 脚本错误原因
	StartTime   time.Time       // 启动时间
	EndTime     time.Time       // 结束时间
}

// 任务执行日志
type JobLog struct {
	JobName      string `bson:"jobName" json:"jobName"`           // 任务名字
	Command      string `bson:"command" json:"command"`           // 脚本命令
	Err          string `bson:"err" json:"err"`                   // 错误原因
	Output       string `bson:"output" json:"output"`             // 脚本输出
	PlanTime     int64  `bson:"planTime" json:"planTime"`         // 计划开始时间，毫秒
	ScheduleTime int64  `bson:"scheduleTime" json:"scheduleTime"` // 实际调度时间
	StartTime    int64  `bson:"startTime" json:"startTime"`       // 任务执行开始时间
	EndTime      int64  `bson:"endTime" json:"endTime"`           // 任务执行结束时间
}
