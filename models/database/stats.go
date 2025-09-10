package database

import (
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

// 访问记录相关方法
func AddVisitRecord(record *VisitRecord) bool {
	db.Create(record)
	return true
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
	// 按session_id去重统计
	query.Group("session_id").Count(&count)
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
	// 按session_id去重统计
	query.Group("session_id").Count(&count)
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
		LanguageStats:       languageStats,
		DeviceStats:         deviceStats,
		CountryStats:        countryStats,
	}, nil
}

func (visitRecord *VisitRecord) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}
