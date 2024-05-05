package start_server

import (
	"database/sql"
	"fmt"

	"github.com/bqdanh/money_transfer/configs/server"
	"github.com/bqdanh/money_transfer/internal/adapters/distribute_lock"
	accountrepo "github.com/bqdanh/money_transfer/internal/adapters/repository/accounts"
	"github.com/bqdanh/money_transfer/internal/adapters/repository/transactions"
	usersrepo "github.com/bqdanh/money_transfer/internal/adapters/repository/users"
	"github.com/bqdanh/money_transfer/internal/adapters/user_token"
	"github.com/bqdanh/money_transfer/internal/adapters/username_pw_validator"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/generate_user_token"
	"github.com/bqdanh/money_transfer/internal/applications/users/validate_username_password"
	"github.com/bqdanh/money_transfer/pkg/database"
	pkgredis "github.com/bqdanh/money_transfer/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type InfrastructureDependencies struct {
	db          *sql.DB
	redisClient *redis.Client
}

func InitInfrastructure(cfg *server.Config) (*InfrastructureDependencies, error) {
	db, err := database.NewMysqlDatabaseConn(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}
	redisClient, err := pkgredis.NewRedisClient(cfg.RedisConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to init redis client: %w", err)
	}

	return &InfrastructureDependencies{
		db:          db,
		redisClient: redisClient,
	}, nil
}

type Adapters struct {
	UserJWT                                 user_token.JWTToken
	UserMysqlRepository                     usersrepo.UserMysqlRepository
	AccountMysqlRepository                  accountrepo.AccountMysqlRepository
	TransactionMysqlRepository              transactions.TransactionMysqlRepository
	DistributeLockWithRedis                 distribute_lock.DistributeLockWithRedis
	ValidateUserNamePasswordWithUserUseCase username_pw_validator.ValidateUserNamePasswordWithUserUseCase
	GenerateUserToken                       generate_user_token.GenerateUserToken
}

func NewAdapters(cfg *server.Config, infra *InfrastructureDependencies) (*Adapters, error) {
	userMysqlRepo, err := usersrepo.NewUserMysqlRepository(infra.db)
	if err != nil {
		return nil, fmt.Errorf("failed to new user repository: %w", err)
	}

	accountMysqlrepo, err := accountrepo.NewAccountMysqlRepository(infra.db)
	if err != nil {
		return nil, fmt.Errorf("failed to new account repository: %w", err)
	}

	jwtTokenAdapter, err := user_token.NewJWTToken(cfg.JwtToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new jwt token: %w", err)
	}
	validateUsernamePasswordHandler, err := validate_username_password.NewValidateUsernamePassword(userMysqlRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to new validate username password application: %w", err)
	}

	usernamePasswordValidatorAdapter, err := username_pw_validator.NewValidateUserNamePasswordWithUserUseCase(validateUsernamePasswordHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to new username password validator: %w", err)
	}

	tokenGenerator, err := generate_user_token.NewGenerateUserToken(jwtTokenAdapter, cfg.GenerateToken)
	if err != nil {
		return nil, fmt.Errorf("failed to new token generator: %w", err)
	}

	distributeLockWithRedis := distribute_lock.NewDistributeLockWithRedis(cfg.DistributeLock, infra.redisClient)

	transactionMysqlRepository, err := transactions.NewTransactionMysqlRepository(infra.db)
	if err != nil {
		return nil, fmt.Errorf("failed to new transaction repository: %w", err)
	}

	return &Adapters{
		UserJWT:                                 jwtTokenAdapter,
		UserMysqlRepository:                     userMysqlRepo,
		AccountMysqlRepository:                  accountMysqlrepo,
		TransactionMysqlRepository:              transactionMysqlRepository,
		DistributeLockWithRedis:                 distributeLockWithRedis,
		ValidateUserNamePasswordWithUserUseCase: usernamePasswordValidatorAdapter,
		GenerateUserToken:                       tokenGenerator,
	}, nil
}
