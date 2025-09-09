package api

import (
	"net/http"
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
	behavior := models.GetUserBehaviorStats()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": behavior,
	})
}
