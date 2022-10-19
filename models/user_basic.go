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
