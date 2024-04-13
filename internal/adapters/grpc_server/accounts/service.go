package accounts

import (
	"github.com/bqdanh/money_transfer/api/grpc/account"
	"github.com/bqdanh/money_transfer/internal/applications/accounts/link_account"
	"google.golang.org/grpc"
)

type AccountService struct {
	account.UnimplementedAccountServiceServer
	App AccountApplications
}

type AccountApplications struct {
	LinkAccount link_account.LinkBankAccount
}

func NewAccountService(app AccountApplications) *AccountService {
	return &AccountService{
		App: app,
	}
}

func (s *AccountService) RegisterService(server grpc.ServiceRegistrar) {
	account.RegisterAccountServiceServer(server, s)
}
