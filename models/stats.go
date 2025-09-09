package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// VisitRecord 访问记录模型
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
	Language  string    `json:"language" gorm:"size:10"`
}

// 访问记录相关方法
func AddVisitRecord(record *VisitRecord) bool {
	db.Create(record)
	return true
}

func GetTodayVisits(language string) int {
	var count int
	today := time.Now().Format("2006-01-02")
	query := db.Model(&VisitRecord{}).Where("DATE(visit_time) = ?", today)
	if language != "" {
		query = query.Where("language = ?", language)
	} else {
		// 当没有指定语言时，只统计有语言信息的记录
		query = query.Where("language IS NOT NULL AND language != ''")
	}
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
	query.Count(&count)
	return count
}

func GetUniqueVisitorsToday(language string) int {
	var count int
	today := time.Now().Format("2006-01-02")
	query := db.Model(&VisitRecord{}).Where("DATE(visit_time) = ?", today)
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
