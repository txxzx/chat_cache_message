package main

import "github.com/txxzx/chat_cache_message/router"

/**
    @date: 2022/10/7
**/

func main() {
	e := router.Router()
	e.Run(":8082")
}
