package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webbleen/go-gin/models/database"
	"github.com/webbleen/go-gin/pkg/e"
)

// GetVisitRecords 获取访问记录列表
// @Summary 获取访问记录列表
// @Description 分页获取访问记录，支持按语言过滤
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param language query string false "语言过滤"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/records [get]
func GetVisitRecords(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	language := c.Query("language")

	// 调用模型层函数
	result, err := database.GetVisitRecords(page, pageSize, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "Failed to get records",
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": result,
	})
}

// GetVisitOverview 获取访问统计概览
// @Summary 获取访问统计概览
// @Description 获取今日访问量、累计访问量、独立访客等统计信息
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Router /stats/overview [get]
func GetVisitOverview(c *gin.Context) {
	// 调用模型层函数
	result, err := database.GetVisitOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "Failed to get overview",
			"data": make(map[string]interface{}),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": result,
	})
}

// DashboardPage 显示 Dashboard 页面
// @Summary 显示 Dashboard 页面
// @Description 显示访问记录管理界面
// @Tags Dashboard
// @Accept html
// @Produce html
// @Success 200 {string} string "HTML页面"
// @Router /dashboard [get]
func DashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "访问记录 Dashboard",
	})
}
