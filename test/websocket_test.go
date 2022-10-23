package test

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

/**
    @date: 2022/10/17
**/
// 服务所在的地址
var addr = flag.String("addr", "localhost:8082", "http service address")

// 通过upgraader 将http封装成一个websocket
var upgrader = websocket.Upgrader{} // use default options

var ws = make(map[*websocket.Conn]struct{})

func TestWebSocketServer(t *testing.T) {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	// 每当一个连接进来了我们就把连接存储进来
	ws[c] = struct{}{}
	for {
		// 接收消息
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// 循环websocket
		for conn := range ws {
			log.Printf("recv: %s", message)
			// 对我们的connection进行发送操作
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

	}
}

func TestGinWebsocketServer(t *testing.T) {
	r := gin.Default()
	// 路由
	r.GET("/echo", func(ctx *gin.Context) {
		echo(ctx.Writer, ctx.Request)
	})
	r.Run(":8083")
}
