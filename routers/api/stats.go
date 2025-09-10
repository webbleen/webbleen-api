package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webbleen/go-gin/models/database"
	"github.com/webbleen/go-gin/pkg/e"
)

// GetVisitStats 获取访问统计概览
// @Summary 获取访问统计概览
// @Description 获取今日访问量、累计访问量、独立访客等统计信息，支持按语言过滤
// @Tags 统计
// @Accept json
// @Produce json
// @Param language query string false "语言代码" default("")
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/visits [get]
func GetVisitStats(c *gin.Context) {
	data := make(map[string]interface{})

	// 获取语言参数
	language := c.Query("language")

	// 今日访问量
	todayVisits := database.GetTodayVisits(language)
	data["today_visits"] = todayVisits

	// 累计访问量
	totalVisits := database.GetTotalVisits(language)
	data["total_visits"] = totalVisits

	// 今日独立访客
	uniqueVisitorsToday := database.GetUniqueVisitorsToday(language)
	data["unique_visitors_today"] = uniqueVisitorsToday

	// 添加语言信息
	if language != "" {
		data["language"] = language
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": data,
	})
}

// RecordVisit 记录访问
// @Summary 记录访问
// @Description 记录用户访问信息，包括页面路径、设备信息、地理位置等
// @Tags 统计
// @Accept json
// @Produce json
// @Param visitRecord body database.VisitRecord true "访问记录"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/visit [post]
func RecordVisit(c *gin.Context) {
	var visitRecord database.VisitRecord

	// 从 JSON 请求体中获取访问信息
	if err := c.ShouldBindJSON(&visitRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid JSON data",
			"data": make(map[string]interface{}),
		})
		return
	}

	// 设置服务器端信息
	// 优先使用前端传递的真实外网IP，其次使用Netlify传递的真实IP，最后使用服务器检测的IP
	if visitRecord.IP == "" {
		// 尝试从Netlify请求头获取真实IP
		if netlifyIP := c.GetHeader("X-Nf-Client-Connection-Ip"); netlifyIP != "" {
			visitRecord.IP = netlifyIP
		} else {
			// 使用Gin的ClientIP()方法，它会自动检查X-Forwarded-For等标准请求头
			visitRecord.IP = c.ClientIP()
		}
	}
	visitRecord.UserAgent = c.GetHeader("User-Agent")
	visitRecord.Referer = c.GetHeader("Referer")

	// 保存访问记录
	database.AddVisitRecord(&visitRecord)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": make(map[string]interface{}),
	})
}

// GetUserBehavior 获取用户行为分析
// @Summary 获取用户行为分析
// @Description 获取设备、浏览器、操作系统、地理位置等用户行为统计
// @Tags 统计
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/behavior [get]
func GetUserBehavior(c *gin.Context) {
	behavior := database.GetUserBehaviorStats()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": behavior,
	})
}
