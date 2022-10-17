package define

import (
	"os"
	"time"
)

/**
    @date: 2022/10/15
**/

var MailPassword = os.Getenv("MailPassword")

var RegisterPrefix = "TOKEN_"
var ExpireTime = time.Second * 300
