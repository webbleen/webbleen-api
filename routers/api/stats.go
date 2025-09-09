package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webbleen/go-gin/models"
	"github.com/webbleen/go-gin/pkg/e"
)

// GetVisitStats 获取访问统计概览
// @Summary 获取访问统计概览
// @Description 获取今日访问量、累计访问量、独立访客等统计信息
// @Tags 统计
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/visits [get]
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

// RecordVisit 记录访问
// @Summary 记录访问
// @Description 记录用户访问信息，包括页面路径、设备信息、地理位置等
// @Tags 统计
// @Accept json
// @Produce json
// @Param visitRecord body models.VisitRecord true "访问记录"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/visit [post]
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

// GetContentStats 获取内容统计
// @Summary 获取内容统计
// @Description 获取博客文章、标签、分类等内容统计信息
// @Tags 统计
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/content [get]
func GetContentStats(c *gin.Context) {
	stats := models.GetContentStats()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": stats,
	})
}

// GetTopPages 获取热门页面
// @Summary 获取热门页面
// @Description 获取访问量最高的页面列表
// @Tags 统计
// @Accept json
// @Produce json
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/pages [get]
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

// GetVisitTrend 获取访问趋势
// @Summary 获取访问趋势
// @Description 获取指定天数的访问趋势数据
// @Tags 统计
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/trend [get]
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

// GetUserBehavior 获取用户行为分析
// @Summary 获取用户行为分析
// @Description 获取设备、浏览器、操作系统、地理位置等用户行为统计
// @Tags 统计
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/behavior [get]
func GetUserBehavior(c *gin.Context) {
	behavior := models.GetUserBehaviorStats()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": behavior,
	})
}

// GetDailyStats 获取日统计
// @Summary 获取日统计
// @Description 获取每日统计数据
// @Tags 统计
// @Accept json
// @Produce json
// @Param limit query int false "限制天数" default(30)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/daily [get]
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

// UpdateContentStats 更新内容统计
// @Summary 更新内容统计
// @Description 用于同步Hugo博客的内容统计信息
// @Tags 统计
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param articles formData int true "文章数量"
// @Param tags formData int true "标签数量"
// @Param categories formData int true "分类数量"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/content [post]
func UpdateContentStats(c *gin.Context) {
	articles := c.PostForm("articles")
	tags := c.PostForm("tags")
	categories := c.PostForm("categories")

	// 转换为整数
	articlesInt, _ := strconv.Atoi(articles)
	tagsInt, _ := strconv.Atoi(tags)
	categoriesInt, _ := strconv.Atoi(categories)

	models.UpdateContentStats(articlesInt, tagsInt, categoriesInt)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": make(map[string]interface{}),
	})
}
