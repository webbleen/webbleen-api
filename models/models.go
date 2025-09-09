package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/webbleen/go-gin/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init() {
	var err error

	// 使用 setting 包中的数据库配置
	db, err = gorm.Open("postgres", setting.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 自动迁移统计相关表
	db.AutoMigrate(&VisitRecord{}, &PageView{}, &DailyStats{})
}

func CloseDB() {
	defer db.Close()
}
