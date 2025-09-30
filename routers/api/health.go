package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/webbleen/go-gin/models/database"
    "github.com/webbleen/go-gin/pkg/e"
)

// Healthz 健康检查（存活）
// @Summary 健康检查
// @Description 返回服务运行状态
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Router /healthz [get]
func Healthz(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "code": e.SUCCESS,
        "msg":  e.GetMsg(e.SUCCESS),
        "data": gin.H{
            "status": "ok",
        },
    })
}

// Readyz 就绪检查（依赖可用）
// @Summary 就绪检查
// @Description 检查数据库等依赖是否可用
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 503 {object} map[string]interface{} "未就绪"
// @Router /readyz [get]
func Readyz(c *gin.Context) {
    // 数据库连通性检查
    db := database.GetDB()
    sqlDB := db.DB()
    if err := sqlDB.Ping(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "code": e.ERROR,
            "msg":  "database not ready",
            "data": gin.H{
                "status": "not_ready",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": e.SUCCESS,
        "msg":  e.GetMsg(e.SUCCESS),
        "data": gin.H{
            "status": "ready",
        },
    })
}

