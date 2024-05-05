package users

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/api/grpc/account"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type AccountService struct {
	grpcServerConn *grpc.ClientConn
}

func NewAccountGatewayService(conn *grpc.ClientConn) *AccountService {
	return &AccountService{
		grpcServerConn: conn,
	}
}

func (s *AccountService) HTTPGatewayRegister(mux *runtime.ServeMux) error {
	err := account.RegisterAccountServiceHandler(context.Background(), mux, s.grpcServerConn)
	if err != nil {
		return fmt.Errorf("failed to register http gateway for account service: %w", err)
	}
	return nil
}
