package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_http_requests_total",
			Help: "Total number of HTTP requests in auth-service",
		},
		[]string{"method", "endpoint"},
	)

	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_http_request_duration_seconds",
			Help:    "HTTP request duration in auth-service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)
