package test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

/**
    @date: 2022/10/16
**/
var ctx = context.Background()

func TestGet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	t.Logf("%s", r)
}

func TestSet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", time.Second*30).Err()
	if err != nil {
		panic(err)
	}
}
