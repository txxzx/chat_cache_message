package service

import (
	"context"
	"github.com/gin-gonic/gin"
	tp "github.com/henrylee2cn/teleport"
	"github.com/txxzx/chat_cache_message/define"
	"github.com/txxzx/chat_cache_message/helper"
	"github.com/txxzx/chat_cache_message/models"
	"net/http"
	"time"
)

/**
    @date: 2022/10/15
**/

// 用户登录的方法
func Login(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	ub, err := models.GetUserBasicByAccountPassword(account, helper.GetMd5(password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}
	// 生成用户token
	token, err := helper.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误:" + err.Error(),
		})
		return
	}
	// 用户登录成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

// 获取用户详情方法
func UserDetail(c *gin.Context) {
	u, _ := c.Get("user_claims")
	uc := u.(*helper.UserClaims)
	userBasic, err := models.GetUserBasicByIdentity(uc.Identity)
	if err != nil {
		tp.Errorf("[DB ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据查询成功",
		"data": userBasic,
	})
}

// 发送验证码方法
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱不能为空",
		})
		return
	}
	count, err := models.GetUserBasicCountByEmail(email)
	if err != nil {
		tp.Errorf("[DB ERROR]:%v", err)
		return
	}
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "当前邮箱已被注册",
		})
		return
	}
	code := helper.GetCode()
	if err := helper.SendCode1(email, code); err != nil {
		tp.Errorf("err-> %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误" + err.Error(),
		})
		return
	}
	if err := models.RedisDB.Set(context.Background(), define.RegisterPrefix+email, code, define.ExpireTime).Err(); err != nil {
		tp.Errorf("[REDIS ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误" + err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
	})
}

// 用户注册方法
func Register(c *gin.Context) {
	code := c.PostForm("code")
	email := c.PostForm("email")
	account := c.PostForm("account")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")
	if code == "" || email == "" || account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	// 验证验证码是否正确
	r, err := models.RedisDB.Get(context.Background(), define.RegisterPrefix+email).Result()
	if err != nil {
		tp.Errorf("[REDIS ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确",
		})
		return
	}
	if r != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确",
		})
		return
	}
	// 判断账号是否正确
	cnt, err := models.GetUserBasicCountByAccount(account)
	if err != nil {
		tp.Errorf("[REDIS ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误" + err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号已被注册",
		})
		return
	}
	ub := &models.UserBasic{
		Identity:  helper.GetUuid(),
		Account:   account,
		Password:  helper.GetMd5(password),
		Nickname:  nickname,
		Email:     email,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	if err := models.InsertOneUserBasic(ub); err != nil {
		tp.Errorf("[REDIS ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误" + err.Error(),
		})
		return
	}
	// 生成用户token
	token, err := helper.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误:" + err.Error(),
		})
		return
	}
	// 用户登录成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"token": token,
		},
	})
}

// 获取用户个人信息
func UserQuery(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询参数不正确",
		})
	}
	ub, err := models.GetUserBasicByAccount(account)
	if err != nil {
		tp.Errorf("[DB ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	uc := c.MustGet("user_claims").(*helper.UserClaims)

	data := UserQueryResult{
		Nickname: ub.Nickname,
		Sex:      ub.Sex,
		Email:    ub.Email,
		Avatar:   ub.Avatar,
		IsFriend: false,
	}
	b, err := models.JudgeUserIsFriend(ub.Identity, uc.Identity)
	if err != nil {
		data.IsFriend = b
	}
	data.IsFriend = b
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据查询成功",
		"data": data,
	})
}

type UserQueryResult struct {
	Nickname string `json:"nickname"`
	Sex      int    `json:"sex"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	IsFriend bool   `json:"is_friend"` // 是否为好友[true 好友 false不是好友]
}

func UserAdd(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	ub, err := models.GetUserBasicByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	uc:= c.MustGet("user_claims").(*helper.UserClaims)

	b, err := models.JudgeUserIsFriend(ub.Identity, uc.Identity)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"数据异常"+err.Error(),
		})
		return
	}
	if b{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"互为好友，不可重复添加",
		})
		return
	}
	// 保存房间记录
	rb := &models.RoomBasic{
		Identity: helper.GetUuid(),
		UserIdentity: uc.Identity,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	if err :=models.InsertOneRoomBasic(rb);err!=nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"数据查询失败",
		})
		return
	}
	// 保存用户与房间的关联关系
	ur := &models.UserRoom{
		UserIdentity: uc.Identity,
		RoomIdentity: rb.Identity,
		RoomType: 1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	if err := models.InsertUserRoom(ur);err!=nil{
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"数据库异常"+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"添加成功",
	})
}
