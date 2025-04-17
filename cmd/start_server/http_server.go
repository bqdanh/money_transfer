package start_server

import (
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/server/http_gateway"
	accountgw "github.com/bqdanh/money_transfer/internal/adapters/server/http_gateway/accounts"
	transactionsgw "github.com/bqdanh/money_transfer/internal/adapters/server/http_gateway/transactions"
	usersgw "github.com/bqdanh/money_transfer/internal/adapters/server/http_gateway/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewHTTPGatewayServices(cfg Config, _ *InfrastructureDependencies) ([]http_gateway.GrpcGatewayServices, error) {
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
	transactionHttpGwService := transactionsgw.NewTransactionGatewayService(grpcServerConn)

	return []http_gateway.GrpcGatewayServices{
		userHttpGwService,
		accountHttpGwService,
		transactionHttpGwService,
	}, nil
}
