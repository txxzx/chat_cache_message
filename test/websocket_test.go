package test

import (
	"flag"
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
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
