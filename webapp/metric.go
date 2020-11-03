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

var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "the number of HTTP requests processed",
	Type:        "counter_vec",
	Args:        []string{"status"}}
