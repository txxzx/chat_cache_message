package service

import (
	"github.com/gin-gonic/gin"
	"github.com/txxzx/chat_cache_message/helper"
	"github.com/txxzx/chat_cache_message/models"
	"net/http"
	"strconv"
)

/**
    @date: 2022/10/23
**/

func ChatList(c *gin.Context) {
	roomIdentity := c.Query("room_identity")
	if roomIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间不能为空",
		})
		return
	}
	// 判断用户是否属于该房间
	// 拿取UserIdentity
	uc := c.MustGet("user_claims").(*helper.UserClaims)
	_, err := models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, roomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "非法访问",
		})
		return
	}
	pageIndex, _ := strconv.ParseInt(c.Query("page_index"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	skip := (pageIndex - 1) * pageSize
	// 聊天记录查询
	data, err := models.GetMessageListByRoomIdentity(roomIdentity, &pageSize, &skip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": gin.H{
			"list": data,
		},
	})
}
