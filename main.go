// @title webbleen 博客 API 服务
// @version 1.0
// @description webbleen 博客 API 服务，提供统计、内容管理等功能
// @termsOfService http://swagger.io/terms/

// @contact.name webbleen
// @contact.url https://webbleen.com
// @contact.email contact@webbleen.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host api.webbleen.com
// @BasePath /
// @schemes https

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
