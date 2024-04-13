package start_server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bqdanh/money_transfer/configs/server"
	grpcadapter "github.com/bqdanh/money_transfer/internal/adapters/grpc_server"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/accounts"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/users"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/utils/authentication_interceptor"
	"github.com/bqdanh/money_transfer/internal/adapters/http_gateway"
	accountgw "github.com/bqdanh/money_transfer/internal/adapters/http_gateway/accounts"
	usersgw "github.com/bqdanh/money_transfer/internal/adapters/http_gateway/users"
	usersrepo "github.com/bqdanh/money_transfer/internal/adapters/repository/users"
	"github.com/bqdanh/money_transfer/internal/adapters/user_token"
	"github.com/bqdanh/money_transfer/internal/adapters/username_pw_validator"
	"github.com/bqdanh/money_transfer/internal/applications/accounts/link_account"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/generate_user_token"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/login"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/validate_user_token"
	"github.com/bqdanh/money_transfer/internal/applications/users/create_user"
	"github.com/bqdanh/money_transfer/internal/applications/users/validate_username_password"
	"github.com/bqdanh/money_transfer/pkg/database"
	"github.com/bqdanh/money_transfer/pkg/logger"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	StartServerCmd = &cli.Command{
		Name:   "server",
		Usage:  "run http server",
		Action: StartServerAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from file path`",
				DefaultText: "./configs/server/local.yaml",
				Value:       "./configs/server/local.yaml",
				Required:    false,
			},
		},
	}
)

type InfrastructureDependencies struct {
	db *sql.DB
}

func StartServerAction(cmdCLI *cli.Context) error {
	cfgPath := cmdCLI.String("config")
	cfg, err := server.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config from path\"%s\": %w", cfgPath, err)
	}
	return StartHTTPServer(cfg)
}

func StartHTTPServer(cfg *server.Config) error {
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

	// new application
	grpcServices, err := NewGrpcServices(*cfg, infra)
	if err != nil {
		return fmt.Errorf("failed to new grpc services: %w", err)
	}

	authenticateInterceptor, err := NewAuthenticateGrpcInterceptors(cfg, infra)
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

	httpStop, cherr := http_gateway.StartServer(cfg.HTTP, httpgwServices...)
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

func NewAuthenticateGrpcInterceptors(cfg *server.Config, _ *InfrastructureDependencies) (grpc.UnaryServerInterceptor, error) {
	jwtTokenAdapter, err := user_token.NewJWTToken(cfg.JwtToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new jwt token: %w", err)
	}

	validator, err := validate_user_token.NewValidateUserToken(jwtTokenAdapter)
	if err != nil {
		return nil, fmt.Errorf("failed to new validate user token: %w", err)
	}
	methodNoNeedAuthenticate := []string{
		"/money_transfer.user_service.UserService/Login",
		"/money_transfer.user_service.UserService/CreateUser",
	}
	authenticateHandler := authentication_interceptor.NewAuthenticationWithUserToken(validator, methodNoNeedAuthenticate)
	return authenticateHandler.UserTokenAuthenticationInterceptor(), nil
}

func InitInfrastructure(cfg *server.Config) (*InfrastructureDependencies, error) {
	db, err := database.NewMysqlDatabaseConn(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	return &InfrastructureDependencies{
		db: db,
	}, nil
}

func NewGrpcServices(cfg server.Config, infra *InfrastructureDependencies) ([]grpcadapter.Service, error) {
	userrepo, err := usersrepo.NewUserMysqlRepository(infra.db)
	if err != nil {
		return nil, fmt.Errorf("failed to new user repository: %w", err)
	}

	// new application
	// user application
	createUserHandler, err := create_user.NewCreateUser(userrepo)
	if err != nil {
		return nil, fmt.Errorf("failed to new create user application: %w", err)
	}
	validateUsernamePasswordHandler, err := validate_username_password.NewValidateUsernamePassword(userrepo)
	if err != nil {
		return nil, fmt.Errorf("failed to new validate username password application: %w", err)
	}

	usernamePasswordValidatorAdapter, err := username_pw_validator.NewValidateUserNamePasswordWithUserUseCase(validateUsernamePasswordHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to new username password validator: %w", err)
	}
	jwtTokenAdapter, err := user_token.NewJWTToken(cfg.JwtToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new jwt token: %w", err)
	}
	tokenGenerator, err := generate_user_token.NewGenerateUserToken(jwtTokenAdapter, cfg.GenerateToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new token generator: %w", err)
	}

	loginHandler, err := login.NewLogin(usernamePasswordValidatorAdapter, tokenGenerator)
	if err != nil {
		return nil, fmt.Errorf("failed to new login application: %w", err)
	}

	userApplications := users.UserApplications{
		CreateUserHandler: createUserHandler,
		Login:             loginHandler,
	}
	// new server
	// new account service
	userService := users.NewUserService(userApplications)

	//account application
	accountApplication := accounts.AccountApplications{
		LinkAccount: link_account.LinkBankAccount{},
	}
	// new account service
	accountService := accounts.NewAccountService(accountApplication)

	return []grpcadapter.Service{
		userService,
		accountService,
	}, nil
}

func NewHTTPGatewayServices(cfg server.Config, _ *InfrastructureDependencies) ([]http_gateway.Services, error) {
	grpcServerAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	grpcServerConn, err := grpc.Dial(grpcServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
		//grpc.WithChainUnaryInterceptor(),
		//grpc.WithUnaryInterceptor(),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to dial gRPC server(%s): %w", grpcServerAddr, err)
	}

	// new http gateway services
	userHttpGwService := usersgw.NewUserGatewayService(grpcServerConn)
	accountHttpGwService := accountgw.NewAccountGatewayService(grpcServerConn)

	return []http_gateway.Services{
		userHttpGwService,
		accountHttpGwService,
	}, nil
}
