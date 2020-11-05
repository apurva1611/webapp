package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

var defaultMetricPath = "/metrics"

// var reqCnt = &Metric{
// 	ID:          "reqCnt",
// 	Name:        "requests_total",
// 	Description: "the number of HTTP requests processed",
// 	Type:        "counter_vec",
// 	Args:        []string{"status"}}

// var (
// 	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
// 		Name:    "example_request_duration_seconds",
// 		Help:    "Histogram for the runtime of a simple example function.",
// 		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
// 	})
// )

var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total_n",
	Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "handler", "host", "url"}}

var reqDur = &Metric{
	ID:          "reqDur",
	Name:        "request_duration_seconds",
	Description: "The HTTP request latencies in seconds.",
	Type:        "histogram_vec",
	Args:        []string{"code", "method", "url"},
}

var resSz = &Metric{
	ID:          "resSz",
	Name:        "response_size_bytes",
	Description: "The HTTP response sizes in bytes.",
	Type:        "summary"}

var reqSz = &Metric{
	ID:          "reqSz",
	Name:        "request_size_bytes",
	Description: "The HTTP request sizes in bytes.",
	Type:        "summary"}

var standardMetrics = []*Metric{
	reqCnt,
	reqDur,
	resSz,
	reqSz,
}
