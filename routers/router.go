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

	// 配置模板引擎
	r.LoadHTMLGlob("web/templates/*")

	// 配置静态文件服务
	r.Static("/static", "./web/static")

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
		// 获取用户行为分析
		stats.GET("/behavior", api.GetUserBehavior)
		// Dashboard API
		stats.GET("/records", api.GetVisitRecords)
		stats.GET("/overview", api.GetVisitOverview)
	}

	// Dashboard 页面
	r.GET("/dashboard", api.DashboardPage)

	return r
}
