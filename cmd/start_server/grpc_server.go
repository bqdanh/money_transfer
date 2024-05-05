package start_server

import (
	"fmt"

	"github.com/bqdanh/money_transfer/configs/server"
	grpcadapter "github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server"
	"github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server/accounts"
	"github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server/transactions"
	"github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server/users"
	"github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server/utils/authentication_interceptor"
	"github.com/bqdanh/money_transfer/internal/applications/accounts/link_account"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/login"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/validate_user_token"
	"github.com/bqdanh/money_transfer/internal/applications/transactions/deposit/make_transaction"
	"github.com/bqdanh/money_transfer/internal/applications/users/create_user"
	"google.golang.org/grpc"
)

func NewGrpcServices(cfg server.Config, infra *InfrastructureDependencies, adapters *Adapters) ([]grpcadapter.Service, error) {
	userService, err := NewUserService(cfg, infra, adapters)
	if err != nil {
		return nil, fmt.Errorf("failed to new user service: %w", err)
	}
	// new account service
	accountService, err := NewAccountService(cfg, infra, adapters)
	if err != nil {
		return nil, fmt.Errorf("failed to new account service: %w", err)
	}
	// new transaction service
	transactionService, err := NewTransactionService(cfg, infra, adapters)

	return []grpcadapter.Service{
		userService,
		accountService,
		transactionService,
	}, nil
}

func NewUserService(_ server.Config, _ *InfrastructureDependencies, adapters *Adapters) (*users.UserService, error) {
	createUserHandler, err := create_user.NewCreateUser(adapters.UserMysqlRepository)
	if err != nil {
		return nil, fmt.Errorf("failed to new create user application: %w", err)
	}

	loginHandler, err := login.NewLogin(adapters.ValidateUserNamePasswordWithUserUseCase, adapters.GenerateUserToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new login application: %w", err)
	}

	userApplications := users.UserApplications{
		CreateUserHandler: createUserHandler,
		Login:             loginHandler,
	}
	userService := users.NewUserService(userApplications)
	return userService, nil
}

func NewAccountService(cfg server.Config, _ *InfrastructureDependencies, adapters *Adapters) (*accounts.AccountService, error) {
	accountApplication := accounts.AccountApplications{
		LinkAccount: link_account.NewLinkBankAccount(cfg.LinkAccount, adapters.AccountMysqlRepository, adapters.DistributeLockWithRedis),
	}
	// new account service
	accountService := accounts.NewAccountService(accountApplication)
	return accountService, nil
}

func NewTransactionService(_ server.Config, _ *InfrastructureDependencies, adapters *Adapters) (*transactions.TransactionService, error) {
	transactionsApp := transactions.TransactionApplications{
		MakeTransactionSync: make_transaction.MakeTransactionSync{},
	}
	// mew transaction service
	transactionService := transactions.NewTransactionService(transactionsApp)
	return transactionService, nil
}

func NewAuthenticateGrpcInterceptors(_ *server.Config, _ *InfrastructureDependencies, adapters *Adapters) (grpc.UnaryServerInterceptor, error) {
	validator, err := validate_user_token.NewValidateUserToken(adapters.UserJWT)
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
