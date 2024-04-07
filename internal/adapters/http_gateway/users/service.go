package users

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/api/grpc/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type UserService struct {
	grpcServerConn *grpc.ClientConn
}

func NewUserGatewayService(conn *grpc.ClientConn) *UserService {
	return &UserService{
		grpcServerConn: conn,
	}
}

func (s *UserService) HTTPGatewayRegister(mux *runtime.ServeMux) error {
	err := user_service.RegisterUserServiceHandler(context.Background(), mux, s.grpcServerConn)
	if err != nil {
		return fmt.Errorf("failed to register http gateway for user service: %w", err)
	}
	return nil
}
