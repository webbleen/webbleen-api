package api

import (
    "encoding/csv"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/webbleen/go-gin/models/database"
)

// ExportVisitRecords 导出访问记录为 CSV
// @Summary 导出访问记录 CSV
// @Description 按页导出访问记录，支持语言过滤
// @Tags 导出
// @Produce text/csv
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(100)
// @Param language query string false "语言过滤"
// @Success 200 {string} string "CSV 文件"
// @Router /stats/export [get]
func ExportVisitRecords(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
    language := c.Query("language")

    result, err := database.GetVisitRecords(page, pageSize, language)
    if err != nil {
        c.String(http.StatusInternalServerError, "failed to query records")
        return
    }

    c.Header("Content-Type", "text/csv; charset=utf-8")
    c.Header("Content-Disposition", "attachment; filename=visit_records.csv")

    writer := csv.NewWriter(c.Writer)
    // 表头
    _ = writer.Write([]string{"id", "ip", "user_agent", "referer", "page", "session_id", "country", "city", "device", "browser", "os", "language", "created_on", "modified_on"})

    for _, r := range result.Records {
        _ = writer.Write([]string{
            strconv.Itoa(r.ID),
            r.IP,
            r.UserAgent,
            r.Referer,
            r.Page,
            r.SessionID,
            r.Country,
            r.City,
            r.Device,
            r.Browser,
            r.OS,
            r.Language,
            r.CreatedOn,
            r.ModifiedOn,
        })
    }

    writer.Flush()
}

