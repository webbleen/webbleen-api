package database

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/webbleen/go-gin/models/response"
)

// VisitRecord 访问记录模型
type VisitRecord struct {
	Model
	IP        string `json:"ip" gorm:"size:45"`
	UserAgent string `json:"user_agent" gorm:"size:500"`
	Referer   string `json:"referer" gorm:"size:500"`
	Page      string `json:"page" gorm:"size:200"`
	SessionID string `json:"session_id" gorm:"size:100"`
	Country   string `json:"country" gorm:"size:50"`
	City      string `json:"city" gorm:"size:50"`
	Device    string `json:"device" gorm:"size:50"`
	Browser   string `json:"browser" gorm:"size:50"`
	OS        string `json:"os" gorm:"size:50"`
	Language  string `json:"language" gorm:"size:10"`
}

// ContentStats 内容统计表
type ContentStats struct {
    Model
    TotalArticles   int       `json:"total_articles"`
    TotalTags       int       `json:"total_tags"`
    TotalCategories int       `json:"total_categories"`
    LastUpdate      time.Time `json:"last_update"`
}

// 访问记录相关方法
func AddVisitRecord(record *VisitRecord) bool {
	// 在存储前解析URL，将编码的路径转换为可读格式
	record.Page = ParseURL(record.Page)
	db.Create(record)
	return true
}

// CheckVisitExists 检查今日是否已记录过该页面的访问
func CheckVisitExists(sessionID, page string) bool {
	today := time.Now().Format("2006-01-02")
	// 解析URL，确保比较的是解析后的格式
	parsedPage := ParseURL(page)
	var count int64
	db.Model(&VisitRecord{}).
		Where("session_id = ? AND page = ? AND DATE(created_on) = ?", sessionID, parsedPage, today).
		Count(&count)
	return count > 0
}

// ParseURL 解析URL，将编码的路径转换为可读格式
func ParseURL(rawURL string) string {
	// 如果URL为空或只是斜杠，返回原值
	if rawURL == "" || rawURL == "/" {
		return rawURL
	}

	// 尝试解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// 如果解析失败，尝试直接解码路径部分
		decoded, decodeErr := url.QueryUnescape(rawURL)
		if decodeErr != nil {
			// 如果解码也失败，返回原值
			return rawURL
		}
		return decoded
	}

	// 解码路径部分
	decodedPath, err := url.QueryUnescape(parsedURL.Path)
	if err != nil {
		// 如果解码失败，返回原路径
		return parsedURL.Path
	}

	// 如果URL有查询参数，保留它们
	if parsedURL.RawQuery != "" {
		decodedQuery, err := url.QueryUnescape(parsedURL.RawQuery)
		if err == nil {
			return decodedPath + "?" + decodedQuery
		}
		return decodedPath + "?" + parsedURL.RawQuery
	}

	return decodedPath
}

func GetTodayVisits(language string) int {
	var count int
	today := time.Now().Format("2006-01-02")
	query := db.Model(&VisitRecord{}).Where("DATE(created_on) = ?", today)
	if language != "" {
		query = query.Where("language = ?", language)
	} else {
		// 当没有指定语言时，只统计有语言信息的记录
		query = query.Where("language IS NOT NULL AND language != ''")
	}
	// 统计所有页面访问（不按session_id去重，每个页面访问都算一次）
	query.Count(&count)
	return count
}

func GetTotalVisits(language string) int {
	var count int
	query := db.Model(&VisitRecord{})
	if language != "" {
		query = query.Where("language = ?", language)
	} else {
		// 当没有指定语言时，只统计有语言信息的记录
		query = query.Where("language IS NOT NULL AND language != ''")
	}
	// 统计所有页面访问（不按session_id去重，每个页面访问都算一次）
	query.Count(&count)
	return count
}

func GetUniqueVisitorsToday(language string) int {
	var count int
	today := time.Now().Format("2006-01-02")
	query := db.Model(&VisitRecord{}).Where("DATE(created_on) = ?", today)
	if language != "" {
		query = query.Where("language = ?", language)
	} else {
		// 当没有指定语言时，只统计有语言信息的记录
		query = query.Where("language IS NOT NULL AND language != ''")
	}
	query.Group("ip").Count(&count)
	return count
}

// GetTodayUniqueSessions 获取今日独立会话数（按session_id去重）
func GetTodayUniqueSessions(language string) int {
	var count int
	today := time.Now().Format("2006-01-02")
	query := db.Model(&VisitRecord{}).Where("DATE(created_on) = ?", today)
	if language != "" {
		query = query.Where("language = ?", language)
	} else {
		// 当没有指定语言时，只统计有语言信息的记录
		query = query.Where("language IS NOT NULL AND language != ''")
	}
	query.Group("session_id").Count(&count)
	return count
}

// GetTotalUniqueSessions 获取总独立会话数（按session_id去重）
func GetTotalUniqueSessions(language string) int {
	var count int
	query := db.Model(&VisitRecord{})
	if language != "" {
		query = query.Where("language = ?", language)
	} else {
		// 当没有指定语言时，只统计有语言信息的记录
		query = query.Where("language IS NOT NULL AND language != ''")
	}
	query.Group("session_id").Count(&count)
	return count
}

// 用户行为分析
func GetUserBehaviorStats() *response.UserBehaviorResult {
	// 设备统计
	var deviceStats []response.DeviceStat
	db.Model(&VisitRecord{}).Select("device, count(*) as count").Group("device").Find(&deviceStats)

	// 浏览器统计
	var browserStats []response.BrowserStat
	db.Model(&VisitRecord{}).Select("browser, count(*) as count").Group("browser").Find(&browserStats)

	// 操作系统统计
	var osStats []response.OSStat
	db.Model(&VisitRecord{}).Select("os, count(*) as count").Group("os").Find(&osStats)

	// 地理位置统计
	var locationStats []response.LocationStat
	db.Model(&VisitRecord{}).Select("country, city, count(*) as count").Group("country, city").Order("count DESC").Limit(10).Find(&locationStats)

	return &response.UserBehaviorResult{
		Devices:          deviceStats,
		Browsers:         browserStats,
		OperatingSystems: osStats,
		Locations:        locationStats,
	}
}

// 热门页面统计（可限制数量）
func GetTopPages(limit int, startDate, endDate string, language string) ([]response.PageStat, error) {
    if limit <= 0 || limit > 100 {
        limit = 10
    }

    query := db.Model(&VisitRecord{})
    if startDate != "" {
        query = query.Where("DATE(created_on) >= ?", startDate)
    }
    if endDate != "" {
        query = query.Where("DATE(created_on) <= ?", endDate)
    }
    if language != "" {
        query = query.Where("language = ?", language)
    }

    type row struct {
        Page  string
        Count int
    }
    var rows []row
    err := query.Select("page, count(*) as count").
        Group("page").
        Order("count DESC").
        Limit(limit).
        Scan(&rows).Error
    if err != nil {
        return nil, err
    }

    stats := make([]response.PageStat, 0, len(rows))
    for _, r := range rows {
        stats = append(stats, response.PageStat{Page: r.Page, Count: r.Count})
    }
    return stats, nil
}

// 趋势/日统计（按天聚合）
func GetTrend(days int, language string) (*response.TrendResult, error) {
    if days <= 0 || days > 365 {
        days = 30
    }
    // 计算起始日期（含当天）
    start := time.Now().AddDate(0, 0, -days+1).Format("2006-01-02")

    type row struct {
        Date  string
        Count int
    }

    // 访问量（不去重）
    query := db.Model(&VisitRecord{}).Where("DATE(created_on) >= ?", start)
    if language != "" {
        query = query.Where("language = ?", language)
    }
    var visitRows []row
    err := query.Select("DATE(created_on) as date, COUNT(*) as count").
        Group("DATE(created_on)").
        Order("date").
        Scan(&visitRows).Error
    if err != nil {
        return nil, err
    }

    // 独立访客（按 IP 去重）
    var uvRows []row
    err = db.Model(&VisitRecord{}).
        Where("DATE(created_on) >= ?", start).
        Scopes(withLanguage(language)).
        Select("DATE(created_on) as date, COUNT(DISTINCT ip) as count").
        Group("DATE(created_on)").
        Order("date").
        Scan(&uvRows).Error
    if err != nil {
        return nil, err
    }

    // 独立会话（按 session_id 去重）
    var usRows []row
    err = db.Model(&VisitRecord{}).
        Where("DATE(created_on) >= ?", start).
        Scopes(withLanguage(language)).
        Select("DATE(created_on) as date, COUNT(DISTINCT session_id) as count").
        Group("DATE(created_on)").
        Order("date").
        Scan(&usRows).Error
    if err != nil {
        return nil, err
    }

    // 合并到完整连续的日期序列
    visitMap := make(map[string]int)
    for _, r := range visitRows { visitMap[r.Date] = r.Count }
    uvMap := make(map[string]int)
    for _, r := range uvRows { uvMap[r.Date] = r.Count }
    usMap := make(map[string]int)
    for _, r := range usRows { usMap[r.Date] = r.Count }

    points := make([]response.TrendPoint, 0, days)
    startTime, _ := time.Parse("2006-01-02", start)
    for i := 0; i < days; i++ {
        d := startTime.AddDate(0, 0, i).Format("2006-01-02")
        points = append(points, response.TrendPoint{
            Date:           d,
            Visits:         visitMap[d],
            UniqueVisitors: uvMap[d],
            UniqueSessions: usMap[d],
        })
    }
    return &response.TrendResult{Points: points}, nil
}

func withLanguage(language string) func(*gorm.DB) *gorm.DB {
    return func(tx *gorm.DB) *gorm.DB {
        if language != "" {
            return tx.Where("language = ?", language)
        }
        return tx
    }
}

// 内容统计读
func GetContentStats() (*response.ContentStatsResponse, error) {
    var cs ContentStats
    // 仅取最新一条
    err := db.Order("modified_on DESC").First(&cs).Error
    if err != nil && !gorm.IsRecordNotFoundError(err) {
        return nil, err
    }
    res := &response.ContentStatsResponse{
        TotalArticles:   cs.TotalArticles,
        TotalTags:       cs.TotalTags,
        TotalCategories: cs.TotalCategories,
    }
    if !cs.ModifiedOn.IsZero() {
        res.LastUpdate = cs.ModifiedOn.Format("2006-01-02 15:04:05")
    }
    return res, nil
}

// 内容统计写（新增一条快照）
func UpdateContentStats(articles, tags, categories int) error {
    cs := ContentStats{
        TotalArticles:   articles,
        TotalTags:       tags,
        TotalCategories: categories,
        LastUpdate:      time.Now(),
    }
    return db.Create(&cs).Error
}
// 分页获取访问记录
func GetVisitRecords(page, pageSize int, language string) (*response.VisitRecordsResult, error) {
	// 限制每页最大数量
	if pageSize > 100 {
		pageSize = 100
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	query := db.Model(&VisitRecord{})

	// 语言过滤
	if language != "" {
		query = query.Where("language = ?", language)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取记录列表
	var records []VisitRecord
	err := query.Order("created_on DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应结构体
	var responseRecords []response.VisitRecord
	for _, record := range records {
		responseRecords = append(responseRecords, response.VisitRecord{
			ID:         record.ID,
			IP:         record.IP,
			UserAgent:  record.UserAgent,
			Referer:    record.Referer,
			Page:       record.Page,
			SessionID:  record.SessionID,
			Country:    record.Country,
			City:       record.City,
			Device:     record.Device,
			Browser:    record.Browser,
			OS:         record.OS,
			Language:   record.Language,
			CreatedOn:  record.CreatedOn.Format("2006-01-02 15:04:05"),
			ModifiedOn: record.ModifiedOn.Format("2006-01-02 15:04:05"),
		})
	}

	// 计算分页信息
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	return &response.VisitRecordsResult{
		Records: responseRecords,
		Pagination: response.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// 获取访问统计概览
func GetVisitOverview() (*response.VisitOverviewResult, error) {
	// 今日访问量
	todayVisits := GetTodayVisits("")

	// 累计访问量
	totalVisits := GetTotalVisits("")

	// 今日独立访客
	uniqueVisitorsToday := GetUniqueVisitorsToday("")

	// 今日独立会话数
	todayUniqueSessions := GetTodayUniqueSessions("")

	// 总独立会话数
	totalUniqueSessions := GetTotalUniqueSessions("")

	// 按语言统计
	languageStats := make(map[string]int64)
	var languages []string
	db.Model(&VisitRecord{}).
		Select("language").
		Group("language").
		Pluck("language", &languages)

	for _, lang := range languages {
		var count int64
		db.Model(&VisitRecord{}).Where("language = ?", lang).Count(&count)
		languageStats[lang] = count
	}

	// 按设备统计
	deviceStats := make(map[string]int64)
	var devices []string
	db.Model(&VisitRecord{}).
		Select("device").
		Group("device").
		Pluck("device", &devices)

	for _, device := range devices {
		var count int64
		db.Model(&VisitRecord{}).Where("device = ?", device).Count(&count)
		deviceStats[device] = count
	}

	// 按国家统计
	countryStats := make(map[string]int64)
	var countries []string
	db.Model(&VisitRecord{}).
		Select("country").
		Group("country").
		Pluck("country", &countries)

	for _, country := range countries {
		var count int64
		db.Model(&VisitRecord{}).Where("country = ?", country).Count(&count)
		countryStats[country] = count
	}

	return &response.VisitOverviewResult{
		TodayVisits:         todayVisits,
		TotalVisits:         totalVisits,
		UniqueVisitorsToday: uniqueVisitorsToday,
		TodayUniqueSessions: todayUniqueSessions,
		TotalUniqueSessions: totalUniqueSessions,
		LanguageStats:       languageStats,
		DeviceStats:         deviceStats,
		CountryStats:        countryStats,
	}, nil
}

func (visitRecord *VisitRecord) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("CreatedOn", now)
	scope.SetColumn("ModifiedOn", now)
	return nil
}
