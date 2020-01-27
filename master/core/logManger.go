package core

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/leiwenxuan/crontab/infra/base"

	"github.com/leiwenxuan/crontab/master/services"
)

var _ services.LoggerMongerServer = new(LogServer)

var once sync.Once

func init() {
	once.Do(func() {
		services.ILoggerMangerServer = new(LogServer)
	})
	logrus.Debug("LogServer 注册", services.ILoggerMangerServer)
}

type LogServer struct {
}

func (l LogServer) ListLog(taskName string, skip int64, limit int64) (jobLog []services.JobLog, err error) {
	client := base.ClientMongodb()

	fmt.Println("client: ", client)
	conf := base.Props()
	database := conf.GetDefault("mongodb.database", "cron")
	collection := conf.GetDefault("mongodb.collection", "log")
	logCollection := client.Database(database).Collection(collection)
	filter := bson.D{{"jobName", taskName}}
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.Skip = &skip

	cur, err := logCollection.Find(context.TODO(), filter, findOptions)
	for cur.Next(context.TODO()) {
		var ruselt services.JobLog
		err := cur.Decode(&ruselt)
		if err != nil {
			log.Fatal(err)
		}
		jobLog = append(jobLog, ruselt)
	}
	return jobLog, err
}
