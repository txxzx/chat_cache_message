package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	tp "github.com/henrylee2cn/teleport"
	"github.com/txxzx/chat_cache_message/define"
	"github.com/txxzx/chat_cache_message/helper"
	"github.com/txxzx/chat_cache_message/models"
	"net/http"
	"time"
)

/**
    @date: 2022/10/22
**/
var upgrader = websocket.Upgrader{}
var wc = make(map[string]*websocket.Conn)

func WebsocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
		return
	}
	defer conn.Close()
	uc := c.MustGet("user_claims").(helper.UserClaims)
	wc[uc.Identity] = conn
	for {
		ms := new(define.MessageStruct)
		err := conn.ReadJSON(ms)
		//  判断用户是否属于消息体的房间
		_, err = models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, ms.RoomIdentity)
		if err != nil {
			tp.Errorf("userIdentity-> %s,roomIdentity-> %s err-> %v", uc.Identity, ms.RoomIdentity, err)
			return
		}
		if err != nil {
			tp.Errorf("%v", err)
			return
		}
		// TODO 保存消息
		mb := &models.MessageBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: ms.RoomIdentity,
			Data:         ms.Message,
			CreatedAt:    time.Now().UnixNano(),
			UpdatedAt:    time.Now().UnixNano(),
		}
		err = models.InsertOneMessageBasic(mb)
		if err != nil {
			tp.Errorf("%v", err)
			return
		}
		// TODO 获取特定房间的在线用户
		room, err := models.GetUserRoomByRoomIdentity(ms.RoomIdentity)
		if err != nil {
			tp.Errorf("room_identity ->%v,err-> %v", ms.RoomIdentity, err)
			return
		}
		for _, r := range room {
			if w, ok := wc[r.UserIdentity]; ok {
				err := w.WriteMessage(websocket.TextMessage, []byte(ms.Message))
				if err != nil {
					tp.Errorf("%v", err)
					return
				}
			}
		}
	}
}
