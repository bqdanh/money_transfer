package http_gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

type Config struct {
	Host string `json:"host" mapstructure:"host" yaml:"host"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}

type Services interface {
	HTTPGatewayRegister(mux *runtime.ServeMux) error
}

func StartServer(cfg Config, services ...Services) (graceShutdown func(), cerr chan error) {
	chanErr := make(chan error, 1)
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(HeaderMatcher),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithMetadata(ExtractInfoAnnotator),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseProtoNames:   true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	for _, service := range services {
		err := service.HTTPGatewayRegister(mux)
		if err != nil {
			chanErr <- fmt.Errorf("failed to register http gateway: %v", err)
			return nil, chanErr
		}
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)

	httpMux, err := RegisterMonitorHttpService(httpMux)
	if err != nil {
		chanErr <- fmt.Errorf("failed to register monitor http service: %v", err)
		return nil, chanErr
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: httpMux,
	}

	go func() {
		defer close(chanErr)
		log.Println("http_gateway server is running on", httpServer.Addr)
		defer log.Println("http_gateway server is stopping")
		if err := httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			chanErr <- fmt.Errorf("failed to serve: %v", err)
			return
		}
	}()

	graceShutdown = func() {
		log.Println("http_gateway server is shutting down")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			chanErr <- fmt.Errorf("failed to shutdown http_gateway: %v", err)
		}
	}

	return graceShutdown, chanErr
}

func RegisterMonitorHttpService(httpMux *http.ServeMux) (*http.ServeMux, error) {
	PrometheusHandler(httpMux)
	PprofHandler(httpMux)
	return httpMux, nil
}

func PrometheusHandler(httpMux *http.ServeMux) {
	httpMux.Handle("/metrics", promhttp.Handler())
}

func PprofHandler(httpMux *http.ServeMux) {
	// Register pprof handlers
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

var _Headers = map[string]struct{}{}

func HeaderMatcher(key string) (string, bool) {
	lowKey := strings.ToLower(key)
	if _, ok := _Headers[lowKey]; ok {
		return lowKey, true
	}
	return "", false
}

func ExtractInfoAnnotator(ctx context.Context, _ *http.Request) metadata.MD {
	md := make(map[string]string)
	if method, ok := runtime.RPCMethod(ctx); ok {
		md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
	}
	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		md["pattern"] = pattern // /v1/example/login
	}
	return metadata.New(md)
}
