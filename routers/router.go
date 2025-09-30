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

    // CORS 設定（可配置）
    r.Use(cors.New(cors.Config{
        AllowOrigins:     setting.CORSAllowedOrigins,
        AllowMethods:     setting.CORSAllowedMethods,
        AllowHeaders:     setting.CORSAllowedHeaders,
        AllowCredentials: setting.CORSCredentials,
        MaxAge:           24 * time.Hour,
    }))

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
