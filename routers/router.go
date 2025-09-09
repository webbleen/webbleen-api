package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	ginswagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/webbleen/go-gin/docs"
	"github.com/webbleen/go-gin/pkg/setting"
	v1 "github.com/webbleen/go-gin/routers/api/v1"
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
	api := r.Group("/api")
	{
		// 记录访问
		api.POST("/visit", v1.RecordVisit)
		// 获取访问统计
		api.GET("/stats/visits", v1.GetVisitStats)
		// 获取内容统计
		api.GET("/stats/content", v1.GetContentStats)
		// 获取热门页面
		api.GET("/stats/pages", v1.GetTopPages)
		// 获取访问趋势
		api.GET("/stats/trend", v1.GetVisitTrend)
		// 获取用户行为分析
		api.GET("/stats/behavior", v1.GetUserBehavior)
		// 获取日统计
		api.GET("/stats/daily", v1.GetDailyStats)
		// 更新内容统计
		api.POST("/stats/content", v1.UpdateContentStats)
	}

	return r
}
