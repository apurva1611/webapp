package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type Prometheus struct {
	reqCnt        *prometheus.CounterVec
	router        *gin.Engine
	listenAddress string

	Metric      *Metric
	MetricsPath string
}

func newPrometheus(subsystem string) *Prometheus {
	p := &Prometheus{
		Metric:        reqCnt,
		MetricsPath:   "/metrics",
		listenAddress: ":9901",
	}

	p.registerMetrics(subsystem)
	p.router = gin.Default()

	return p
}

func (p *Prometheus) registerMetrics(subsystem string) {
	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      reqCnt.Name,
			Help:      reqCnt.Description,
		},
		reqCnt.Args,
	)
	if err := prometheus.Register(metric); err != nil {
		log.Infof("%s could not be registered: ", reqCnt, err)
	} else {
		log.Infof("%s registered.", reqCnt)
	}
	p.reqCnt = metric

	reqCnt.MetricCollector = metric
}
