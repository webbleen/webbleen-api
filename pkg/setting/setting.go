package setting

import (
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

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
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadDatabase()
    LoadCORS()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func LoadDatabase() {
	// 从环境变量读取数据库连接
	DatabaseURL = os.Getenv("DATABASE_URL")
	if DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
}

func LoadCORS() {
    sec, err := Cfg.GetSection("cors")
    if err != nil {
        // 未配置 cors 段落则使用默认值
        CORSAllowedOrigins = []string{
            "http://127.0.0.1:8080",
            "http://localhost:1313",
        }
        CORSAllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
        CORSAllowedHeaders = []string{"Content-Type", "Authorization", "Accept-Encoding", "Content-Length"}
        CORSCredentials = true
        return
    }

    CORSAllowedOrigins = splitAndTrim(sec.Key("ALLOWED_ORIGINS").MustString("http://127.0.0.1:8080,http://localhost:1313"))
    CORSAllowedMethods = splitAndTrim(sec.Key("ALLOWED_METHODS").MustString("GET,POST,PUT,DELETE,OPTIONS"))
    CORSAllowedHeaders = splitAndTrim(sec.Key("ALLOWED_HEADERS").MustString("Content-Type,Authorization,Accept-Encoding,Content-Length"))
    CORSCredentials = sec.Key("ALLOW_CREDENTIALS").MustBool(true)
}

func splitAndTrim(s string) []string {
    var result []string
    cur := ""
    for i := 0; i < len(s); i++ {
        if s[i] == ',' {
            if len(cur) > 0 {
                result = append(result, trimSpaces(cur))
                cur = ""
            }
        } else {
            cur += string(s[i])
        }
    }
    if len(cur) > 0 {
        result = append(result, trimSpaces(cur))
    }
    return result
}

func trimSpaces(s string) string {
    start := 0
    end := len(s)
    for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
        start++
    }
    for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
        end--
    }
    return s[start:end]
}
