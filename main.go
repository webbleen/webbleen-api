package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/webbleen/go-gin/pkg/setting"
	"github.com/webbleen/go-gin/routers"
)

func main() {
	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20

	// 优先使用 Railway 的 PORT 环境变量
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(setting.HTTPPort)
	}

	endPoint := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on port %s", port)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
