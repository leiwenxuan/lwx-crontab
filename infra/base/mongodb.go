package base

import (
	"context"

	"github.com/leiwenxuan/crontab/infra"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb日志管理
var clientMongo *mongo.Client

func ClientMongodb() *mongo.Client {
	Check(clientMongo)
	return clientMongo
}

type MongoDBStarter struct {
	infra.BaseStarter
}

func (s *MongoDBStarter) Setup(ctx infra.StarterContext) {
	var err error
	conf := ctx.Props()
	applyURI := conf.GetDefault("mongodb.applyURI", "mongodb://root:123456@192.168.1.5:10000")
	clientOptions := options.Client().ApplyURI(applyURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Error("mongo 连接失败", err)
		return
	}
	clientMongo = client
}
