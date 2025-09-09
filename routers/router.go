package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	ginswagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/webbleen/go-gin/docs"
	"github.com/webbleen/go-gin/pkg/setting"
	"github.com/webbleen/go-gin/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// Cors設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1:8080",
			"http://192.168.0.7:8080",
			"http://localhost:1313",
			"https://webbleen.com",
			"https://www.webbleen.com",
			"https://webbleen.github.io",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerFiles.Handler))

	// 统计相关API - 不需要认证
	stats := r.Group("/stats")
	{
		// 记录访问
		stats.POST("/visit", api.RecordVisit)
		// 获取访问统计
		stats.GET("/visits", api.GetVisitStats)
		// 获取热门页面
		stats.GET("/pages", api.GetTopPages)
		// 获取访问趋势
		stats.GET("/trend", api.GetVisitTrend)
		// 获取用户行为分析
		stats.GET("/behavior", api.GetUserBehavior)
		// 获取日统计
		stats.GET("/daily", api.GetDailyStats)
	}

	return r
}
