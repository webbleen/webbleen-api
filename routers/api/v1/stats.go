package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webbleen/go-gin/models"
	"github.com/webbleen/go-gin/pkg/e"
)

// 获取访问统计概览
func GetVisitStats(c *gin.Context) {
	data := make(map[string]interface{})

	// 今日访问量
	todayVisits := models.GetTodayVisits()
	data["today_visits"] = todayVisits

	// 累计访问量
	totalVisits := models.GetTotalVisits()
	data["total_visits"] = totalVisits

	// 今日独立访客
	uniqueVisitorsToday := models.GetUniqueVisitorsToday()
	data["unique_visitors_today"] = uniqueVisitorsToday

	// 页面浏览量
	pageViews := models.GetPageViews()
	data["page_views"] = pageViews

	// 总页面浏览量
	var totalPageViews int
	for _, pv := range pageViews {
		totalPageViews += pv.ViewCount
	}
	data["total_page_views"] = totalPageViews

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": data,
	})
}

// 记录访问
func RecordVisit(c *gin.Context) {
	var visitRecord models.VisitRecord

	// 从请求中获取访问信息
	visitRecord.IP = c.ClientIP()
	visitRecord.UserAgent = c.GetHeader("User-Agent")
	visitRecord.Referer = c.GetHeader("Referer")
	visitRecord.Page = c.Query("page")
	visitRecord.SessionID = c.Query("session_id")
	visitRecord.Country = c.Query("country")
	visitRecord.City = c.Query("city")
	visitRecord.Device = c.Query("device")
	visitRecord.Browser = c.Query("browser")
	visitRecord.OS = c.Query("os")
	visitRecord.VisitTime = time.Now()

	// 保存访问记录
	models.AddVisitRecord(&visitRecord)

	// 更新页面访问统计
	if visitRecord.Page != "" {
		models.IncrementPageView(visitRecord.Page)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": make(map[string]interface{}),
	})
}

// 获取内容统计
func GetContentStats(c *gin.Context) {
	stats := models.GetContentStats()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": stats,
	})
}

// 获取热门页面
func GetTopPages(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		limit = int(c.GetInt("limit"))
	}

	pages := models.GetTopPages(limit)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": pages,
	})
}

// 获取访问趋势
func GetVisitTrend(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		days = int(c.GetInt("days"))
	}

	trend := models.GetVisitTrend(days)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": trend,
	})
}

// 获取用户行为分析
func GetUserBehavior(c *gin.Context) {
	behavior := models.GetUserBehaviorStats()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": behavior,
	})
}

// 获取日统计
func GetDailyStats(c *gin.Context) {
	limit := 30
	if l := c.Query("limit"); l != "" {
		limit = int(c.GetInt("limit"))
	}

	stats := models.GetDailyStats(limit)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": stats,
	})
}

// 更新内容统计（用于同步Hugo博客的内容）
func UpdateContentStats(c *gin.Context) {
	articles := c.GetInt("articles")
	tags := c.GetInt("tags")
	categories := c.GetInt("categories")

	models.UpdateContentStats(articles, tags, categories)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": make(map[string]interface{}),
	})
}
