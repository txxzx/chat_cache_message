package helper

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"time"

	"crypto/tls"

	"crypto/md5"
	"math/rand"
	"net/smtp"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
    @date: 2022/10/15
**/

type UserClaims struct {
	Identity string `json:"identity"`
	//Identity primitive.ObjectID `json:"identity"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("im")

// GenerateToken
// 生成 token
func GenerateToken(identity, email string) (string, error) {
	//objectId, err := primitive.ObjectIDFromHex(identity)
	//if err != nil {
	//	return "", err
	//}
	UserClaim := &UserClaims{
		Identity:       identity,
		Email:          email,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// SendCode
// 发送验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <1164165957@qq.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "1164165957@qq.com", "bglkejfqywgmidji", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

// GetCode
// 生成验证码
func GetCode() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// GetUuid
// 生成用户唯一码
func GetUuid() string {
	return primitive.NewObjectID().String()
}

// 发送验证码
func SendCode1(toUserEmail, code string) error {

	sender := "1164165957@qq.com"  //发送者qq邮箱
	authCode := "bglkejfqywgmidji" //qq邮箱授权码
	mailTitle := "IM验证码"           //邮件标题
	mailBody := "你的验证码为" + code    //邮件内容,可以是html

	//接收者邮箱列表
	mailTo := []string{
		toUserEmail,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)       //发送者腾讯邮箱账号
	m.SetHeader("To", mailTo...)      //接收者邮箱列表
	m.SetHeader("Subject", mailTitle) //邮件标题
	m.SetBody("text/html", mailBody)  //邮件内容,可以是html

	// //添加附件
	// zipPath := "./user/temp.zip"
	// m.Attach(zipPath)

	//发送邮件服务器、端口、发送者qq邮箱、qq邮箱授权码
	//服务器地址和端口是腾讯的
	d := gomail.NewDialer("smtp.qq.com", 587, sender, authCode)
	err := d.DialAndSend(m)
	return err
}
