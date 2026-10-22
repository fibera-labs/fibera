package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request processing duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed.",
		},
		[]string{"method"},
	)
)

// Metrics Prometheus 监控中间件，采集请求量、耗时与并发数
func Metrics() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 指标端点自身不参与采集，避免自监控数据干扰
		if c.Path() == "/metrics" {
			return c.Next()
		}

		method := c.Method()
		httpRequestsInFlight.WithLabelValues(method).Inc()
		defer httpRequestsInFlight.WithLabelValues(method).Dec()

		start := time.Now()
		err := c.Next()

		// 使用路由模板作为标签，避免动态路径导致标签基数膨胀
		path := c.Route().Path
		if path == "" {
			path = "unmatched"
		}
		status := strconv.Itoa(c.Response().StatusCode())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path, status).Observe(time.Since(start).Seconds())
		return err
	}
}
