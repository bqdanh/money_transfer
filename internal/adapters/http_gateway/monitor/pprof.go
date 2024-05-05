package monitor

import (
	"net/http"
	"net/http/pprof"
)

type PprofService struct {
}

func (p PprofService) HTTPRegister(httpMux *http.ServeMux) error {
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return nil
}
