package main

import (
	"github.com/gin-gonic/gin"
	"kitexTwo/client/router"
)

func main() {
	r := gin.Default()
	router.Router(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
