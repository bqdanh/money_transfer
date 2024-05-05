package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusService struct {
}

func (p PrometheusService) HTTPRegister(httpMux *http.ServeMux) error {
	httpMux.Handle("/metrics", promhttp.Handler())
	return nil
}
