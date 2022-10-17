package models

import (
	"context"
	"github.com/go-redis/redis/v8"
	tp "github.com/henrylee2cn/teleport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

/**
    @date: 2022/10/15
**/

var Mongo = InitMongo()
var RedisDB = InitRedis()

func InitMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		tp.Errorf("%v", err)
		return nil
	}
	return client.Database("im")
}

// 初始化redis
func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
