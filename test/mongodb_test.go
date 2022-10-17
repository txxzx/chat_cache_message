package test

import (
	"context"
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
