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

	// 获取用户详情
	auth := r.Group("/u", middlewares.AuthCheck())
	auth.GET("/user/detail", service.UserDetail)
	// 发送验证码
	r.POST("/send/code", service.SendCode)
	// 用户注册
	r.POST("/register", service.Register)
	// 查询用户的个人信息
	r.GET("/user/query", service.UserQuery)
	return r
}
