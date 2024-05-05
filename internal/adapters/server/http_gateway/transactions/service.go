package transactions

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/api/grpc/transaction"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type TransactionService struct {
	grpcServerConn *grpc.ClientConn
}

func NewTransactionGatewayService(conn *grpc.ClientConn) *TransactionService {
	return &TransactionService{
		grpcServerConn: conn,
	}
}

func (s *TransactionService) HTTPGatewayRegister(mux *runtime.ServeMux) error {
	err := transaction.RegisterTransactionServiceHandler(context.Background(), mux, s.grpcServerConn)
	if err != nil {
		return fmt.Errorf("failed to register http gateway for transaction service: %w", err)
	}
	return nil
}
