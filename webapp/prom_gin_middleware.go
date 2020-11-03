package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// use adds the middleware to a gin engine.
func (p *Prometheus) use(e *gin.Engine) {
	e.Use(p.handlerFunc())
	p.setMetricsPath(e)
}

func (p *Prometheus) handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() == p.MetricsPath {
			c.Next()
			return
		}
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		p.reqCnt.WithLabelValues(status).Inc()
	}
}

func (p *Prometheus) setMetricsPath(e *gin.Engine) {
	p.router.GET(p.MetricsPath, prometheusHandler())
	go p.router.Run(p.listenAddress)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
