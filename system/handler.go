package system

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct{}

func (metrics) Handler() (string, http.Handler) {
	return "GET /metrics", promhttp.Handler()
}
