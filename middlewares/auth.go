package middlewares

import (
	"github.com/gin-gonic/gin"
	tp "github.com/henrylee2cn/teleport"
	"github.com/txxzx/chat_cache_message/helper"
	"net/http"
)

/**
    @date: 2022/10/15
**/

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		userClaims, err := helper.AnalyseToken(token)
		if err != nil {
			tp.Errorf("err->%v", err)
			// 中断请求
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户认证不通过",
			})
			return
		}
		// 用户认证通过设置用户
		c.Set("user_claims", userClaims)
		c.Next()
	}
}
