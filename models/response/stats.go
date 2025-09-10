package response

// 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// 访问记录分页结果
type VisitRecordsResult struct {
	Records    []VisitRecord `json:"records"`
	Pagination Pagination    `json:"pagination"`
}

// 访问记录（用于响应）
type VisitRecord struct {
	ID         int    `json:"id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Referer    string `json:"referer"`
	Page       string `json:"page"`
	SessionID  string `json:"session_id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Device     string `json:"device"`
	Browser    string `json:"browser"`
	OS         string `json:"os"`
	Language   string `json:"language"`
	CreatedOn  string `json:"created_on"`
	ModifiedOn string `json:"modified_on"`
}

// 访问统计概览结果
type VisitOverviewResult struct {
	TodayVisits         int              `json:"today_visits"`
	TotalVisits         int              `json:"total_visits"`
	UniqueVisitorsToday int              `json:"unique_visitors_today"`
	LanguageStats       map[string]int64 `json:"language_stats"`
	DeviceStats         map[string]int64 `json:"device_stats"`
	CountryStats        map[string]int64 `json:"country_stats"`
}

// 用户行为分析结果
type UserBehaviorResult struct {
	Devices          []DeviceStat   `json:"devices"`
	Browsers         []BrowserStat  `json:"browsers"`
	OperatingSystems []OSStat       `json:"operating_systems"`
	Locations        []LocationStat `json:"locations"`
}

// 设备统计
type DeviceStat struct {
	Device string `json:"device"`
	Count  int    `json:"count"`
}

// 浏览器统计
type BrowserStat struct {
	Browser string `json:"browser"`
	Count   int    `json:"count"`
}

// 操作系统统计
type OSStat struct {
	OS    string `json:"os"`
	Count int    `json:"count"`
}

// 地理位置统计
type LocationStat struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Count   int    `json:"count"`
}
