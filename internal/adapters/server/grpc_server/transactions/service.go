package transactions

import (
	"context"
	"errors"

	"github.com/bqdanh/money_transfer/api/grpc/transaction"
	"github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server/utils"
	"github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server/utils/exceptions_parser"
	"github.com/bqdanh/money_transfer/internal/applications/transactions/deposit/make_transaction"
	"github.com/bqdanh/money_transfer/internal/entities/currency"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	transactionentity "github.com/bqdanh/money_transfer/internal/entities/transaction"
	"github.com/bqdanh/money_transfer/pkg/logger"
	"google.golang.org/grpc"
)

type TransactionService struct {
	transaction.UnimplementedTransactionServiceServer
	App TransactionApplications
}

type TransactionApplications struct {
	make_transaction.MakeTransactionSync
}

func NewTransactionService(app TransactionApplications) *TransactionService {
	return &TransactionService{
		App: app,
	}
}

func (s TransactionService) RegisterService(server grpc.ServiceRegistrar) {
	transaction.RegisterTransactionServiceServer(server, s)
}

func (s TransactionService) MakeDeposit(ctx context.Context, req *transaction.MakeDepositRequest) (*transaction.MakeDepositResponse, error) {
	l := logger.FromContext(ctx)

	currencyUnit, err := currency.GetCurrencyUnit(req.GetCurrency())
	if err != nil {
		rerr := exceptions.NewInvalidArgumentError(
			"currency",
			"invalid currency unit",
			map[string]interface{}{
				"currency": req.GetCurrency(),
				"error":    err,
			},
		)
		return nil, exceptions_parser.Err2GrpcStatus(rerr).Err()
	}
	amount, err := currency.FromFloat64(req.GetAmount(), currencyUnit)
	if err != nil {
		rerr := exceptions.NewInvalidArgumentError(
			"amount",
			"invalid amount",
			map[string]interface{}{
				"amount": req.GetAmount(),
				"error":  err,
			},
		)
		return nil, exceptions_parser.Err2GrpcStatus(rerr).Err()
	}

	trans, err := s.App.MakeTransactionSync.Handle(ctx, make_transaction.MakeDepositTransactionParams{
		RequestID:    req.GetRequestId(),
		UserID:       req.GetUserId(),
		AccountID:    req.GetAccountId(),
		Amount:       amount,
		Descriptions: req.GetDescriptions(),
		Source:       "client",
	})
	if err != nil {
		l.Errorw("failed to create deposit transaction", "error", err)
		if errors.Is(err, make_transaction.ErrFailedToCreateDepositTransaction) {
			return nil, exceptions_parser.Err2GrpcStatus(err).Err()
		}
		if errors.Is(err, make_transaction.ErrFailedToProcessDepositTransaction) {
			if trans.IsProcessing() {
				return &transaction.MakeDepositResponse{
					Code:    utils.CodeSuccess,
					Message: utils.MessageSuccess,
					Data: &transaction.MakeDepositResponse_Data{
						TransactionId:     trans.ID,
						TransactionStatus: fromTransactionStatus2GrpcStatus(trans.Status),
					},
				}, nil
			}
			if trans.IsSuccess() {
				return &transaction.MakeDepositResponse{
					Code:    utils.CodeSuccess,
					Message: utils.MessageSuccess,
					Data: &transaction.MakeDepositResponse_Data{
						TransactionId:     trans.ID,
						TransactionStatus: fromTransactionStatus2GrpcStatus(trans.Status),
					},
				}, nil
			}
			return nil, exceptions_parser.Err2GrpcStatus(err).Err()
		}
		return nil, exceptions_parser.Err2GrpcStatus(err).Err()
	}
	
	return &transaction.MakeDepositResponse{
		Code:    utils.CodeSuccess,
		Message: utils.MessageSuccess,
		Data: &transaction.MakeDepositResponse_Data{
			TransactionId:     trans.ID,
			TransactionStatus: fromTransactionStatus2GrpcStatus(trans.Status),
		},
	}, nil
}

func fromTransactionStatus2GrpcStatus(status transactionentity.Status) transaction.TransactionStatus {
	switch status {
	case transactionentity.StatusInit:
		return transaction.TransactionStatus_PROCESSING
	case transactionentity.StatusProcessing:
		return transaction.TransactionStatus_PROCESSING
	case transactionentity.StatusSuccess:
		return transaction.TransactionStatus_SUCCESS
	case transactionentity.StatusFailed:
		return transaction.TransactionStatus_FAILED
	default:
		return transaction.TransactionStatus_PROCESSING
	}
}
