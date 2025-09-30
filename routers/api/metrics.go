package api

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
    "time"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

// PrometheusHandler 暴露 /metrics
func PrometheusHandler(c *gin.Context) {
    promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

// MetricsMiddleware 记录基础 HTTP 指标
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()

        duration := time.Since(start).Seconds()
        path := c.FullPath()
        if path == "" {
            path = c.Request.URL.Path
        }
        status := c.Writer.Status()

        httpRequestsTotal.WithLabelValues(c.Request.Method, path, http.StatusText(status)).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
    }
}

