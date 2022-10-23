package models

import (
	"context"
	tp "github.com/henrylee2cn/teleport"
	"go.mongodb.org/mongo-driver/bson"
)

/**
    @date: 2022/10/22
**/

type RoomBasic struct {
	Identity     string `bson:"identity"`
	Number       string `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	UserIdentity string `bson:"user_identity"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (RoomBasic) CollectionName() string {
	return "room_basic"
}

// 插入房间信息
func InsertOneRoomBasic(rb *RoomBasic) error {
	_, err := Mongo.Collection(RoomBasic{}.CollectionName()).InsertOne(context.Background(), rb)
	if err != nil {
		return err
	}
	return nil
}

// 删除RoomBasic里面的关联关系
func DeleteRoomBasic(roomIdentity string) error {
	if _, err := Mongo.Collection(RoomBasic{}.CollectionName()).
		DeleteOne(context.Background(), bson.M{"identity": roomIdentity}); err != nil {
		tp.Errorf("%v", err)
		return err
	}
	return nil
}
