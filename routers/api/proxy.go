package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webbleen/go-gin/pkg/e"
)

// BingResponse 必应壁纸响应结构
type BingResponse struct {
	Images []struct {
		URL       string `json:"url"`
		Copyright string `json:"copyright"`
		Title     string `json:"title"`
	} `json:"images"`
}

// FaviconResponse 网站图标响应结构
type FaviconResponse struct {
	URL string `json:"url"`
}

// GeoResponse 地理位置响应结构
type GeoResponse struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Region  string `json:"region"`
	IP      string `json:"ip"`
}

// IPResponse IP地址响应结构
type IPResponse struct {
	IP string `json:"ip"`
}

// GetBingWallpaper 获取必应每日壁纸
// @Summary 获取必应每日壁纸
// @Description 代理访问必应每日壁纸数据
// @Tags 代理服务
// @Accept json
// @Produce json
// @Success 200 {object} e.Response{data=BingResponse}
// @Router /proxy/bing [get]
func GetBingWallpaper(c *gin.Context) {
	// 必应壁纸API
	url := "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "获取必应壁纸失败",
			"data": nil,
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "读取响应失败",
			"data": nil,
		})
		return
	}

	var bingResp BingResponse
	if err := json.Unmarshal(body, &bingResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  "解析响应失败",
			"data": nil,
		})
		return
	}

	// 处理图片URL
	if len(bingResp.Images) > 0 {
		image := bingResp.Images[0]
		if !strings.HasPrefix(image.URL, "http") {
			image.URL = "https://www.bing.com" + image.URL
		}
		bingResp.Images[0] = image
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": bingResp,
	})
}

// GetFavicon 获取网站图标
// @Summary 获取网站图标
// @Description 代理访问网站图标服务
// @Tags 代理服务
// @Accept json
// @Produce json
// @Param url query string true "网站URL"
// @Success 200 {object} e.Response{data=FaviconResponse}
// @Router /proxy/favicon [get]
func GetFavicon(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "缺少URL参数",
			"data": nil,
		})
		return
	}

	// 尝试多个favicon服务
	faviconServices := []string{
		fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=32", url),
		fmt.Sprintf("https://t3.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=%s&size=32", url),
		fmt.Sprintf("https://favicons.githubusercontent.com/%s", url),
	}

	var faviconURL string
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 尝试每个服务直到找到一个可用的
	for _, serviceURL := range faviconServices {
		resp, err := client.Head(serviceURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			faviconURL = serviceURL
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	if faviconURL == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code": e.ERROR,
			"msg":  "无法获取网站图标",
			"data": nil,
		})
		return
	}

	faviconResp := FaviconResponse{
		URL: faviconURL,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": faviconResp,
	})
}

// GetGeoLocation 获取地理位置信息
// @Summary 获取地理位置信息
// @Description 返回客户端的地理位置信息
// @Tags 代理服务
// @Accept json
// @Produce json
// @Success 200 {object} e.Response{data=GeoResponse}
// @Router /proxy/geo [get]
func GetGeoLocation(c *gin.Context) {
	// 获取客户端IP
	clientIP := c.ClientIP()

	// 如果是从Docker内部网络获取的IP，尝试从外部API获取真实IP
	if strings.HasPrefix(clientIP, "192.168.") || strings.HasPrefix(clientIP, "172.") || strings.HasPrefix(clientIP, "10.") || strings.HasPrefix(clientIP, "100.64.") {
		// 使用ipify.org获取真实公网IP
		url := "https://api.ipify.org?format=text"
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil && len(body) > 0 {
				clientIP = strings.TrimSpace(string(body))
			}
		}
	}

	// 尝试多个地理位置服务
	geoServices := []struct {
		name string
		url  string
	}{
		{"ip-api.com", fmt.Sprintf("https://ip-api.com/json/%s", clientIP)},
		{"ipinfo.io", fmt.Sprintf("https://ipinfo.io/%s/json", clientIP)},
		{"ipapi.co", fmt.Sprintf("https://ipapi.co/%s/json/", clientIP)},
	}

	var geoResp GeoResponse

	for _, service := range geoServices {
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(service.url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		var geoData map[string]interface{}
		if err := json.Unmarshal(body, &geoData); err != nil {
			continue
		}

		// 检查是否有错误（不同服务的错误字段不同）
		if _, exists := geoData["error"]; exists {
			continue
		}
		if status, exists := geoData["status"]; exists && status != "success" {
			continue
		}

		// 根据不同的服务解析数据
		switch service.name {
		case "ipapi.co":
			geoResp = GeoResponse{
				Country: getString(geoData, "country_name"),
				City:    getString(geoData, "city"),
				Region:  getString(geoData, "region"),
				IP:      clientIP,
			}
		case "ip-api.com":
			geoResp = GeoResponse{
				Country: getString(geoData, "country"),
				City:    getString(geoData, "city"),
				Region:  getString(geoData, "regionName"),
				IP:      clientIP,
			}
		case "ipinfo.io":
			geoResp = GeoResponse{
				Country: getString(geoData, "country"),
				City:    getString(geoData, "city"),
				Region:  getString(geoData, "region"),
				IP:      clientIP,
			}
		}

		// 如果获取到了有效数据，跳出循环
		if geoResp.Country != "" {
			break
		}
	}

	// 如果所有服务都失败了，返回默认值
	if geoResp.Country == "" {
		geoResp = GeoResponse{
			Country: "Unknown",
			City:    "Unknown",
			Region:  "Unknown",
			IP:      clientIP,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": geoResp,
	})
}

// GetClientIP 获取客户端IP地址
// @Summary 获取客户端IP地址
// @Description 返回客户端的公网IP地址
// @Tags 代理服务
// @Accept json
// @Produce json
// @Success 200 {object} e.Response{data=IPResponse}
// @Router /proxy/ip [get]
func GetClientIP(c *gin.Context) {
	// 尝试从多个来源获取真实IP
	clientIP := c.ClientIP()

	// 如果是从Docker内部网络获取的IP，尝试从外部API获取真实IP
	if strings.HasPrefix(clientIP, "192.168.") || strings.HasPrefix(clientIP, "172.") || strings.HasPrefix(clientIP, "10.") {
		// 使用ipify.org获取真实公网IP
		url := "https://api.ipify.org?format=text"
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil && len(body) > 0 {
				clientIP = strings.TrimSpace(string(body))
			}
		}
	}

	ipResp := IPResponse{
		IP: clientIP,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": ipResp,
	})
}

// 辅助函数：从map中安全获取字符串值
func getString(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}
