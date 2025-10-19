package api

import (
	"net/http"
	"strconv"

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

	// 今日独立会话数
	todayUniqueSessions := database.GetTodayUniqueSessions(language)
	data["today_unique_sessions"] = todayUniqueSessions

	// 总独立会话数
	totalUniqueSessions := database.GetTotalUniqueSessions(language)
	data["total_unique_sessions"] = totalUniqueSessions

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

	// 检查今日是否已记录过该页面的访问
	if database.CheckVisitExists(visitRecord.SessionID, visitRecord.Page) {
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  "Visit already recorded today",
			"data": map[string]interface{}{
				"recorded": false,
				"reason":   "already_exists",
			},
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
		"data": map[string]interface{}{
			"recorded": true,
			"reason":   "new_visit",
		},
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

// GetTopPages 获取热门页面
// @Summary 获取热门页面
// @Description 按访问量降序返回热门页面列表
// @Tags 统计
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Param language query string false "语言过滤"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/pages [get]
func GetTopPages(c *gin.Context) {
	limit := 10
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	start := c.Query("start_date")
	end := c.Query("end_date")
	language := c.Query("language")

	stats, err := database.GetTopPages(limit, start, end, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": e.ERROR, "msg": "Failed to get pages", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": e.SUCCESS, "msg": e.GetMsg(e.SUCCESS), "data": gin.H{"pages": stats}})
}

// GetTrend 获取访问趋势
// @Summary 获取访问趋势
// @Description 返回最近N天的访问量、独立访客、独立会话趋势
// @Tags 统计
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Param language query string false "语言过滤"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/trend [get]
func GetTrend(c *gin.Context) {
	days := 30
	if v := c.Query("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			days = n
		}
	}
	language := c.Query("language")
	res, err := database.GetTrend(days, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": e.ERROR, "msg": "Failed to get trend", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": e.SUCCESS, "msg": e.GetMsg(e.SUCCESS), "data": res})
}

// GetDaily 获取日统计（与趋势同结构，主要用于固定时间窗口或分页）
// @Summary 获取日统计
// @Description 返回按天聚合的访问/访客/会话数据
// @Tags 统计
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Param language query string false "语言过滤"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/daily [get]
func GetDaily(c *gin.Context) {
	days := 30
	if v := c.Query("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			days = n
		}
	}
	language := c.Query("language")
	res, err := database.GetTrend(days, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": e.ERROR, "msg": "Failed to get daily", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": e.SUCCESS, "msg": e.GetMsg(e.SUCCESS), "data": gin.H{"points": res.Points}})
}

// GetContentStats 获取内容统计
// @Summary 获取内容统计
// @Description 返回文章/标签/分类等汇总
// @Tags 统计
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/content [get]
func GetContentStats(c *gin.Context) {
	res, err := database.GetContentStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": e.ERROR, "msg": "Failed to get content stats", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": e.SUCCESS, "msg": e.GetMsg(e.SUCCESS), "data": res})
}

// UpdateContentStats 更新内容统计
// @Summary 更新内容统计
// @Description 更新文章/标签/分类数量（追加一条快照）
// @Tags 统计
// @Accept json
// @Produce json
// @Param body body struct{Articles int `json:"articles"`; Tags int `json:"tags"`; Categories int `json:"categories"`} true "内容统计"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/content [post]
func UpdateContentStats(c *gin.Context) {
	var payload struct {
		Articles   int `json:"articles"`
		Tags       int `json:"tags"`
		Categories int `json:"categories"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Invalid JSON data", "data": gin.H{}})
		return
	}
	if err := database.UpdateContentStats(payload.Articles, payload.Tags, payload.Categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": e.ERROR, "msg": "Failed to update content stats", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": e.SUCCESS, "msg": e.GetMsg(e.SUCCESS), "data": gin.H{"updated": true}})
}
