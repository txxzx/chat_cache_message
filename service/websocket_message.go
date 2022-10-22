package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	tp "github.com/henrylee2cn/teleport"
	"github.com/txxzx/chat_cache_message/define"
	"github.com/txxzx/chat_cache_message/helper"
	"net/http"
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
		if err != nil {
			tp.Errorf("%v", err)
			return
		}
		for _, w := range wc {
			err := w.WriteMessage(websocket.TextMessage, []byte(ms.Message))
			if err != nil {
				tp.Errorf("%v", err)
				return
			}
		}
	}
}
