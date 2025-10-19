package setting

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// 服务器配置
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// 应用配置
	PageSize  int
	JwtSecret string

	// 数据库配置
	DatabaseURL string

	// CORS 配置
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string
	CORSCredentials    bool
)

func init() {
	LoadAll()
}

// LoadAll 加载所有配置
func LoadAll() {
	LoadServer()
	LoadApp()
	LoadDatabase()
	LoadCORS()
}

// LoadServer 加载服务器配置
func LoadServer() {
	// 运行模式
	RunMode = getEnv("GIN_MODE", "debug")

	// 端口配置
	HTTPPort = getEnvInt("PORT", 8000)

	// 超时配置
	ReadTimeout = time.Duration(getEnvInt("READ_TIMEOUT", 60)) * time.Second
	WriteTimeout = time.Duration(getEnvInt("WRITE_TIMEOUT", 60)) * time.Second
}

// LoadApp 加载应用配置
func LoadApp() {
	// JWT 密钥
	JwtSecret = getEnv("JWT_SECRET", "!@)*#)!@U#@*!@!)")

	// 分页大小
	PageSize = getEnvInt("PAGE_SIZE", 10)
}

// LoadDatabase 加载数据库配置
func LoadDatabase() {
	// 数据库连接字符串
	DatabaseURL = getEnv("DATABASE_URL", "")
}

// LoadCORS 加载 CORS 配置
func LoadCORS() {
	// 允许的来源
	origins := getEnv("CORS_ALLOWED_ORIGINS", "")
	CORSAllowedOrigins = splitAndTrim(origins)

	// 允许的方法
	methods := getEnv("CORS_ALLOWED_METHODS", "")
	CORSAllowedMethods = splitAndTrim(methods)

	// 允许的头部
	headers := getEnv("CORS_ALLOWED_HEADERS", "")
	CORSAllowedHeaders = splitAndTrim(headers)

	// 是否允许携带凭据
	CORSCredentials = getEnvBool("CORS_CREDENTIALS", false)
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool 获取环境变量并转换为布尔值，如果不存在或转换失败则返回默认值
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// splitAndTrim 分割字符串并去除空白字符
func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	parts := strings.Split(s, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// PrintConfig 打印当前配置（用于调试）
func PrintConfig() {
	log.Printf("=== 当前配置 ===")
	log.Printf("运行模式: %s", RunMode)
	log.Printf("HTTP 端口: %d", HTTPPort)
	log.Printf("读取超时: %v", ReadTimeout)
	log.Printf("写入超时: %v", WriteTimeout)
	log.Printf("分页大小: %d", PageSize)
	log.Printf("数据库 URL: %s", maskSensitiveInfo(DatabaseURL))
	log.Printf("CORS 允许来源: %v", CORSAllowedOrigins)
	log.Printf("CORS 允许方法: %v", CORSAllowedMethods)
	log.Printf("CORS 允许头部: %v", CORSAllowedHeaders)
	log.Printf("CORS 允许凭据: %t", CORSCredentials)
	log.Printf("================")
}

// maskSensitiveInfo 遮蔽敏感信息
func maskSensitiveInfo(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}
