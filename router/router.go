package router

import (
	"github.com/gin-gonic/gin"
	"github.com/txxzx/chat_cache_message/middlewares"
	"github.com/txxzx/chat_cache_message/service"
)

/**
    @date: 2022/10/15
**/

func Router() *gin.Engine {
	r := gin.Default()
	// 定义一个POST的请求 用户的登录
	r.POST("/login", service.Login)
	// 发送验证码
	r.POST("/send/code", service.SendCode)
	// 用户注册
	r.POST("/register", service.Register)

	// 获取用户详情
	auth := r.Group("/u", middlewares.AuthCheck())
	auth.GET("/user/detail", service.UserDetail)
	// 发送接收消息
	auth.GET("/websocket/message", service.WebsocketMessage)

	// 查询用户的个人信息
	auth.GET("/user/query", service.UserQuery)
	// 聊天记录查询
	auth.GET("/chat/list", service.ChatList)

	// 添加用户
	auth.POST("/user/add", service.UserAdd)

	// 删除好友
	auth.DELETE("/user/delete", service.DeleteUser)
	return r
}
