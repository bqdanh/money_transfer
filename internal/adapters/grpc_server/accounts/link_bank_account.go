package accounts

import (
	"context"

	"github.com/bqdanh/money_transfer/api/grpc/account"
	"github.com/bqdanh/money_transfer/internal/adapters/grpc_server/utils"
	"github.com/bqdanh/money_transfer/internal/applications/accounts/link_account"
	"github.com/bqdanh/money_transfer/pkg/logger"
)

func (s *AccountService) LinkBankAccount(ctx context.Context, req *account.LinkBankAccountRequest) (*account.LinkBankAccountResponse, error) {
	ac, err := s.App.LinkAccount.Handle(ctx, link_account.LinkBankAccountParams{
		UserID:            req.GetUserId(),
		BankCode:          req.GetBankCode(),
		BankAccountNumber: req.GetAccountNumber(),
		BankAccountName:   req.GetAccountName(),
	})
	if err != nil {
		logger.FromContext(ctx).Errorw("link bank account", "err", err)
		return nil, err
	}
	return &account.LinkBankAccountResponse{
		Code:    utils.CodeSuccess,
		Message: utils.MessageSuccess,
		Data: &account.LinkBankAccountResponse_Data{
			AccountId: ac.ID,
		},
	}, nil

}
