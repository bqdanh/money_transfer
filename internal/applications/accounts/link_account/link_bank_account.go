package link_account

import (
	"context"
	"fmt"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	bank_account2 "github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

type Config struct {
	LockDuration time.Duration `json:"lock_duration" mapstructure:"lock_duration"`
}

var DefaultConfig = Config{
	LockDuration: 30 * time.Second,
}

type LinkBankAccount struct {
	cfg               Config
	accountRepository accountRepository
	distributeLock    distributeLock
}

//go:generate mockgen -destination=mocks.go -package=link_account -source=link_bank_account.go .

type accountRepository interface {
	GetAccountsByUserID(ctx context.Context, userID int64) ([]account.Account, error)
	//CreateAccount create account for user, return account with ID, ID is unique
	CreateAccount(ctx context.Context, a account.Account) (account.Account, error)
}

type distributeLock interface {
	AcquireLockForCreateAccountByUserID(ctx context.Context, userID int64, lockDuration time.Duration) (releaseLock func(), err error)
}

func NewLinkBankAccount(cfg Config, a accountRepository, dl distributeLock) LinkBankAccount {
	return LinkBankAccount{
		cfg:               cfg,
		accountRepository: a,
		distributeLock:    dl,
	}
}

type LinkBankAccountParams struct {
	UserID            int64
	BankCode          string
	BankAccountNumber string
	BankAccountName   string
}

func validateLinkBankAccountParams(p LinkBankAccountParams) error {
	if p.UserID <= 0 {
		return exceptions.NewInvalidArgumentError(
			"UserID",
			"user id must be greater than 0",
			map[string]interface{}{
				"UserID": p.UserID,
			},
		)
	}
	if p.BankCode == "" {
		return exceptions.NewInvalidArgumentError(
			"BankCode",
			"bank code must not be empty",
			map[string]interface{}{
				"BankCode": p.BankCode,
			},
		)
	}
	if p.BankAccountNumber == "" {
		return exceptions.NewInvalidArgumentError(
			"BankAccountNumber",
			"bank account number must not be empty",
			map[string]interface{}{
				"BankAccountNumber": p.BankAccountNumber,
			},
		)
	}
	if p.BankAccountName == "" {
		return exceptions.NewInvalidArgumentError(
			"BankAccountName",
			"bank account name must not be empty",
			map[string]interface{}{
				"BankAccountName": p.BankAccountName,
			},
		)
	}
	return nil
}

//TODO: implement logic check user is valid for linking bank account

func (l LinkBankAccount) Handle(ctx context.Context, p LinkBankAccountParams) (account.Account, error) {
	if err := validateLinkBankAccountParams(p); err != nil {
		return account.Account{}, err
	}
	bankSofCode, err := bank_account2.FromStringToSourceOfFundCode(p.BankCode)
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to convert bank code to source of fund code: %w", err)
	}
	newBankSof, err := bank_account2.CreateSourceOfFundBankAccount(bankSofCode, bank_account2.BankAccount{
		AccountNumber: p.BankAccountNumber,
		AccountName:   p.BankAccountName,
	})
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to create source of fund bank account: %w", err)
	}
	releaseLock, err := l.distributeLock.AcquireLockForCreateAccountByUserID(ctx, p.UserID, l.cfg.LockDuration)
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to acquire lock for create account: %w", err)
	}
	defer releaseLock()
	accounts, err := l.accountRepository.GetAccountsByUserID(ctx, p.UserID)
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to get accounts by user id: %w", err)
	}
	for _, a := range accounts {
		if a.Status == account.StatusUnlinked {
			//TODO: currently we relink bank account if it is unlinked, we can change this logic to enable the account is unlinked
			continue
		}
		if a.SourceOfFundData.IsTheSameSof(newBankSof) {
			//TODO: we can add more logic to check status of sof, example: sof is locked, inactive, etc
			return account.Account{}, exceptions.NewPreconditionError(
				exceptions.PreconditionReasonAccountIsLinked,
				exceptions.SubjectAccount,
				"account is already linked",
				map[string]interface{}{
					"UserID":    p.UserID,
					"AccountID": a.ID,
				},
			)
		}
	}

	return l.accountRepository.CreateAccount(ctx, account.Account{
		UserID:           p.UserID,
		Status:           account.StatusNormal,
		SourceOfFundData: newBankSof,
	})
}
