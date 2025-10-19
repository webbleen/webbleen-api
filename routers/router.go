package routers

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	ginswagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/webbleen/go-gin/models/database"
	"github.com/webbleen/go-gin/pkg/setting"
	"github.com/webbleen/go-gin/routers/api"
)

func InitRouter() *gin.Engine {
	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		// 如果数据库初始化失败，记录错误但不中断服务
		// 在实际生产环境中，应该处理这个错误
		log.Printf("数据库初始化失败: %v", err)
	}

	r := gin.New()

	// 配置模板引擎
	r.LoadHTMLGlob("web/templates/*")

	// 配置静态文件服务
	r.Static("/static", "./web/static")

	// CORS 設定（可配置）
	// 只有当配置了 CORS 时才启用 CORS 中间件
	if len(setting.CORSAllowedOrigins) > 0 {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     setting.CORSAllowedOrigins,
			AllowMethods:     setting.CORSAllowedMethods,
			AllowHeaders:     setting.CORSAllowedHeaders,
			AllowCredentials: setting.CORSCredentials,
			MaxAge:           24 * time.Hour,
		}))
	}

	r.Use(gin.Logger())
	r.Use(api.MetricsMiddleware())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/metrics", api.PrometheusHandler)

	// 健康与就绪检查
	r.GET("/healthz", api.Healthz)
	r.GET("/readyz", api.Readyz)

	// 统计相关API - 不需要认证
	stats := r.Group("/stats")
	{
		// 记录访问
		stats.POST("/visit", api.RecordVisit)
		// 获取访问统计
		stats.GET("/visits", api.GetVisitStats)
		// 获取用户行为分析
		stats.GET("/behavior", api.GetUserBehavior)
		// 获取热门页面
		stats.GET("/pages", api.GetTopPages)
		// 获取访问趋势 & 日统计
		stats.GET("/trend", api.GetTrend)
		stats.GET("/daily", api.GetDaily)
		// 内容统计读写
		stats.GET("/content", api.GetContentStats)
		stats.POST("/content", api.UpdateContentStats)
		// Dashboard API
		stats.GET("/records", api.GetVisitRecords)
		stats.GET("/overview", api.GetVisitOverview)
		// 导出 CSV
		stats.GET("/export", api.ExportVisitRecords)
	}

	// 代理服务API - 不需要认证
	proxy := r.Group("/proxy")
	{
		// 必应壁纸
		proxy.GET("/bing", api.GetBingWallpaper)
		// 网站图标
		proxy.GET("/favicon", api.GetFavicon)
		// 地理位置
		proxy.GET("/geo", api.GetGeoLocation)
		// IP地址
		proxy.GET("/ip", api.GetClientIP)
	}

	// Dashboard 页面
	r.GET("/dashboard", api.DashboardPage)

	return r
}
