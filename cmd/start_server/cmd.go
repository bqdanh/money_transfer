package start_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	grpcadapter "github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server"
	"github.com/bqdanh/money_transfer/internal/adapters/server/http_gateway"
	"github.com/bqdanh/money_transfer/internal/adapters/server/http_gateway/monitor"
	"github.com/bqdanh/money_transfer/pkg/logger"
	"github.com/urfave/cli/v2"
)

var (
	Cmd = &cli.Command{
		Name:   "server",
		Usage:  "run http server",
		Action: StartServerAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from file path`",
				DefaultText: "./cmd/start_server/config/local.yaml",
				Value:       "./cmd/start_server/config/local.yaml",
				Required:    false,
			},
		},
	}
)

func StartServerAction(cmdCLI *cli.Context) error {
	cfgPath := cmdCLI.String("config")
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config from path\"%s\": %w", cfgPath, err)
	}
	return StartHTTPServer(cfg)
}

func StartHTTPServer(cfg *Config) error {
	if cfg.Env == "local" {
		bs, err := json.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		fmt.Println("Start server with config:", string(bs))
	}

	err := logger.InitLogger(&cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	l := logger.FromContext(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cerr := make(chan error)
	go func() {
		for _err := range cerr {
			l.Errorw("server error got error", "error", _err)
			stop <- syscall.SIGTERM
		}
	}()

	//init infrastructure
	infra, err := InitInfrastructure(cfg)
	if err != nil {
		return fmt.Errorf("failed to init infrastructure: %w", err)
	}

	adapters, err := NewAdapters(cfg, infra)
	if err != nil {
		return fmt.Errorf("failed to new adapters: %w", err)
	}

	// new application
	grpcServices, err := NewGrpcServices(*cfg, infra, adapters)
	if err != nil {
		return fmt.Errorf("failed to new grpc services: %w", err)
	}

	authenticateInterceptor, err := NewAuthenticateGrpcInterceptors(cfg, infra, adapters)
	if err != nil {
		return fmt.Errorf("failed to new authenticate grpc interceptor: %w", err)
	}
	// start server
	grpcStop, cgrpcerr := grpcadapter.StartServer(cfg.GRPC, authenticateInterceptor, grpcServices...)
	go func() {
		for gerr := range cgrpcerr {
			cerr <- fmt.Errorf("grpc server error: %w", gerr)
		}
	}()
	defer grpcStop()

	// start http server
	httpgwServices, err := NewHTTPGatewayServices(*cfg, infra)
	if err != nil {
		return fmt.Errorf("failed to new http gateway services: %w", err)
	}

	httpServices := []http_gateway.HTTPService{
		monitor.PprofService{},
		monitor.PrometheusService{},
	}
	httpStop, cherr := http_gateway.StartServer(cfg.HTTP, httpgwServices, httpServices)
	go func() {
		for herr := range cherr {
			cerr <- fmt.Errorf("http server error: %w", herr)
		}
	}()
	defer httpStop()

	l.Infow("server started")
	<-stop
	l.Infow("server stopping")
	return nil
}
