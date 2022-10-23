package models

import (
	"context"
	tp "github.com/henrylee2cn/teleport"
	"go.mongodb.org/mongo-driver/bson"
)

/**
    @date: 2022/10/22
**/

type UserRoom struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType int `bson:"room_type"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

// 通过用户标识房间标识查询房间信息
func GetUserRoomByUserIdentityRoomIdentity(UserIdentity, RoomIdentity string) (*UserRoom, error) {
	ub := new(UserRoom)
	if err := Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"user_identity", UserIdentity}, {"room_identity", RoomIdentity}}).
		Decode(ub); err != nil {
		tp.Errorf("%v", err)
		return nil, err
	}
	return ub, nil
}

func GetUserRoomByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	// bson里面添加的查询条件
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{{"room_identity", roomIdentity}})
	if err != nil {
		return nil, err
	}
	us := make([]*UserRoom, 0)
	for cursor.Next(context.Background()) {
		ub := new(UserRoom)
		err := cursor.Decode(ub)
		if err != nil {
			tp.Errorf("%v", err)
		}
		us = append(us, ub)
	}
	return us, nil
}

// 插入房间信息
func InsertUserRoom(ur *UserRoom) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).InsertOne(context.Background(), ur)
	if err != nil {
		return err
	}
	return nil
}
