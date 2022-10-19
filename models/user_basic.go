package models

import (
	"context"
	tp "github.com/henrylee2cn/teleport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
    @date: 2022/10/15
**/

type UserBasic struct {
	Identity  string `bson:"_id"`
	Account   string `bson:"account"`
	Password  string `bson:"password"`
	Nickname  string `bson:"nickname"`
	Sex       int    `bson:"sex"`
	Email     string `bson:"email"`
	Avatar    string `bson:"avatar"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

func (UserBasic) CollectionName() string {
	return "user_basic"
}

// 通过账号和密码获取用户信息
func GetUserBasicByAccountPassword(account, password string) (*UserBasic, error) {
	ub := new(UserBasic)
	if err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"account", account}, {"password", password}}).
		Decode(ub); err != nil {
		tp.Errorf("%v", err)
		return nil, err
	}
	return ub, nil
}

// 通过用户标识查询用户信息
func GetUserBasicByIdentity(identity primitive.ObjectID) (*UserBasic, error) {
	ub := new(UserBasic)
	if err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"_id", identity}}).
		Decode(ub); err != nil {
		tp.Errorf("%v", err)
		return nil, err
	}
	return ub, nil
}

// 通过邮箱查询用户的个数
func GetUserBasicCountByEmail(email string) (int64, error) {
	return Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{{"email", email}})
}

// 通过邮箱查询用户的个数
func GetUserBasicCountByAccount(account string) (int64, error) {
	return Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{{"account", account}})
}

// 插入用户注册信息
func InsertOneUserBasic(ub *UserBasic) error {
	_, err := Mongo.Collection(UserBasic{}.CollectionName()).InsertOne(context.Background(), ub)
	if err != nil {
		return err
	}
	return nil
}

// 根据用户账号查询用户信息
func GetUserBasicByAccount(account string) (*UserBasic, error) {
	ub := new(UserBasic)
	if err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"_id", account}}).
		Decode(ub); err != nil {
		tp.Errorf("%v", err)
		return nil, err
	}
	return ub, nil
}

// 判断用户是否为好友
func JudgeUserIsFriend(UserIdentity1, UserIdentity2 string) (bool, error) {
	// 查询user1的单聊房间列表
	cuser, err := Mongo.Collection(UserRoom{}.CollectionNameRoom()).
		Find(context.Background(), bson.D{{"user_identity", UserIdentity1}, {"room_type", 1}})
	roomIdentity := make([]string, 0)
	if err != nil {
		return false, err
	}
	for cuser.Next(context.Background()) {
		ur := new(UserRoom)
		err := cuser.Decode(ur)
		if err != nil {
			return false, err
		}
		roomIdentity = append(roomIdentity, ur.RoomIdentity)
	}
	// 获取关联 userIdentity2 单间聊天房间数
	cunt, err := Mongo.Collection(UserRoom{}.CollectionNameRoom()).CountDocuments(context.Background(), bson.M{"user_identity": UserIdentity2, "room_identity": bson.M{"$in": roomIdentity}})
	if err != nil {
		return false, err
	}
	if cunt > 0 {
		return true, nil
	}
	return false, nil
}
