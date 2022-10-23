package models

import (
	"context"
	tp "github.com/henrylee2cn/teleport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
    @date: 2022/10/22
**/

type MessageBasic struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (MessageBasic) CollectionName() string {
	return "message_basic"
}

// 通过用户标识房间标识查询房间信息
func InsertOneMessageBasic(mb *MessageBasic) error {
	if _, err := Mongo.Collection(MessageBasic{}.CollectionName()).
		InsertOne(context.Background(), mb); err != nil {
		tp.Errorf("%v", err)
		return err
	}
	return nil
}

func GetMessageListByRoomIdentity(roomIdentity string, limit *int64, skip *int64) ([]*MessageBasic, error) {
	data := make([]*MessageBasic, 0)
	cus, err := Mongo.Collection(MessageBasic{}.CollectionName()).
		// 进行分页查询，并且进行排序
		Find(context.Background(), bson.M{"room_identity": roomIdentity}, &options.FindOptions{
			Limit: limit,
			Skip:  skip,
			Sort: bson.D{{
				"created_at", -1,
			}},
		})
	if err != nil {
		return nil, err
	}
	for cus.Next(context.Background()) {
		mb := new(MessageBasic)
		err := cus.Decode(mb)
		if err != nil {
			return nil, err
		}
		data = append(data, mb)
	}
	return data, nil
}
