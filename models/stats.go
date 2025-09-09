package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 访问记录模型
type VisitRecord struct {
	Model
	IP        string    `json:"ip" gorm:"size:45"`
	UserAgent string    `json:"user_agent" gorm:"size:500"`
	Referer   string    `json:"referer" gorm:"size:500"`
	Page      string    `json:"page" gorm:"size:200"`
	VisitTime time.Time `json:"visit_time"`
	SessionID string    `json:"session_id" gorm:"size:100"`
	Country   string    `json:"country" gorm:"size:50"`
	City      string    `json:"city" gorm:"size:50"`
	Device    string    `json:"device" gorm:"size:50"`
	Browser   string    `json:"browser" gorm:"size:50"`
	OS        string    `json:"os" gorm:"size:50"`
}

// 页面访问统计模型
type PageView struct {
	Model
	Page        string    `json:"page" gorm:"size:200;unique_index"`
	ViewCount   int       `json:"view_count"`
	LastView    time.Time `json:"last_view"`
	UniqueViews int       `json:"unique_views"`
}

// 日访问统计模型
type DailyStats struct {
	Model
	Date           time.Time `json:"date" gorm:"unique_index"`
	PageViews      int       `json:"page_views"`
	UniqueVisitors int       `json:"unique_visitors"`
	BounceRate     float64   `json:"bounce_rate"`
	AvgTime        int       `json:"avg_time"` // 平均停留时间（秒）
}

// 内容统计模型
type ContentStats struct {
	Model
	TotalArticles   int       `json:"total_articles"`
	TotalTags       int       `json:"total_tags"`
	TotalCategories int       `json:"total_categories"`
	LastUpdate      time.Time `json:"last_update"`
}

// 访问记录相关方法
func AddVisitRecord(record *VisitRecord) bool {
	db.Create(record)
	return true
}

func GetTodayVisits() int {
	var count int
	today := time.Now().Format("2006-01-02")
	db.Model(&VisitRecord{}).Where("DATE(visit_time) = ?", today).Count(&count)
	return count
}

func GetTotalVisits() int {
	var count int
	db.Model(&VisitRecord{}).Count(&count)
	return count
}

func GetUniqueVisitorsToday() int {
	var count int
	today := time.Now().Format("2006-01-02")
	db.Model(&VisitRecord{}).Where("DATE(visit_time) = ?", today).Group("ip").Count(&count)
	return count
}

func GetPageViews() []PageView {
	var pageViews []PageView
	db.Order("view_count DESC").Find(&pageViews)
	return pageViews
}

func IncrementPageView(page string) {
	var pageView PageView
	if db.Where("page = ?", page).First(&pageView).RecordNotFound() {
		pageView = PageView{
			Page:        page,
			ViewCount:   1,
			LastView:    time.Now(),
			UniqueViews: 1,
		}
		db.Create(&pageView)
	} else {
		db.Model(&pageView).Updates(map[string]interface{}{
			"view_count": pageView.ViewCount + 1,
			"last_view":  time.Now(),
		})
	}
}

func GetDailyStats(limit int) []DailyStats {
	var stats []DailyStats
	db.Order("date DESC").Limit(limit).Find(&stats)
	return stats
}

func GetContentStats() ContentStats {
	var stats ContentStats
	db.Last(&stats)
	return stats
}

func UpdateContentStats(articles, tags, categories int) {
	stats := ContentStats{
		TotalArticles:   articles,
		TotalTags:       tags,
		TotalCategories: categories,
		LastUpdate:      time.Now(),
	}
	db.Create(&stats)
}

func GetTopPages(limit int) []PageView {
	var pages []PageView
	db.Order("view_count DESC").Limit(limit).Find(&pages)
	return pages
}

func GetVisitTrend(days int) []DailyStats {
	var stats []DailyStats
	startDate := time.Now().AddDate(0, 0, -days)
	db.Where("date >= ?", startDate).Order("date ASC").Find(&stats)
	return stats
}

// 用户行为分析
func GetUserBehaviorStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// 设备统计
	var deviceStats []struct {
		Device string
		Count  int
	}
	db.Model(&VisitRecord{}).Select("device, count(*) as count").Group("device").Find(&deviceStats)
	stats["devices"] = deviceStats

	// 浏览器统计
	var browserStats []struct {
		Browser string
		Count   int
	}
	db.Model(&VisitRecord{}).Select("browser, count(*) as count").Group("browser").Find(&browserStats)
	stats["browsers"] = browserStats

	// 操作系统统计
	var osStats []struct {
		OS    string
		Count int
	}
	db.Model(&VisitRecord{}).Select("os, count(*) as count").Group("os").Find(&osStats)
	stats["operating_systems"] = osStats

	// 地理位置统计
	var locationStats []struct {
		Country string
		City    string
		Count   int
	}
	db.Model(&VisitRecord{}).Select("country, city, count(*) as count").Group("country, city").Order("count DESC").Limit(10).Find(&locationStats)
	stats["locations"] = locationStats

	return stats
}

func (visitRecord *VisitRecord) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (pageView *PageView) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (dailyStats *DailyStats) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (contentStats *ContentStats) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}
