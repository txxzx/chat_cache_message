package test

import (
	"context"
	tp "github.com/henrylee2cn/teleport"
	"github.com/txxzx/chat_cache_message/models"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
    @date: 2022/10/15
**/

func TestFindOne(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Errorf("%v", err)
	}
	db := client.Database("im")
	ub := new(models.UserBasic)
	err = db.Collection("user_basic").FindOne(context.Background(), bson.D{}).Decode(ub)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("ub-%v", ub)

}

func TestFing(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Errorf("%v", err)
	}
	db := client.Database("im")

	cursor, err := db.Collection("user_room").Find(context.Background(), bson.D{})
	us := make([]*models.UserRoom, 0)
	for cursor.Next(context.Background()) {
		ub := new(models.UserRoom)
		err := cursor.Decode(ub)
		if err != nil {
			tp.Errorf("%v", err)
		}
		us = append(us, ub)
	}
	for _, v := range us {
		tp.Infof("%v", v)
	}

}
