package link_account

import (
	"context"
	"fmt"
	"testing"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account/implement_bank_account"
	_ "github.com/bqdanh/money_transfer/internal/entities/sof/bank_account/implement_bank_account"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLinkBankAccount_Handle(t *testing.T) {
	uterr := fmt.Errorf("ut force error")
	type args struct {
		ctx context.Context
		p   LinkBankAccountParams
	}
	type want struct {
		ac  account.Account
		err error
	}
	type depmock struct {
		accountRepository func(tt *testing.T) *MockaccountRepository
		distributeLock    func(tt *testing.T) *MockdistributeLock
	}
	tests := []struct {
		name    string
		args    args
		depmock depmock
		want    want
	}{
		{
			name: "validate params UserID error",
			args: args{
				p: LinkBankAccountParams{},
			},
			want: want{
				ac: account.Account{},
				err: exceptions.NewInvalidArgumentError(
					"UserID",
					"user id must be greater than 0",
					map[string]interface{}{
						"UserID": 0,
					},
				),
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
			},
		},
		{
			name: "validate params BankCode error",
			args: args{
				p: LinkBankAccountParams{
					UserID: 1,
				},
			},
			want: want{
				ac: account.Account{},
				err: exceptions.NewInvalidArgumentError(
					"BankCode",
					"bank code must not be empty",
					map[string]interface{}{
						"BankCode": "",
					},
				),
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
			},
		},
		{
			name: "validate params BankAccountNumber error",
			args: args{
				p: LinkBankAccountParams{
					UserID:   1,
					BankCode: "VIB",
				},
			},
			want: want{
				ac: account.Account{},
				err: exceptions.NewInvalidArgumentError(
					"BankAccountNumber",
					"bank account number must not be empty",
					map[string]interface{}{
						"BankAccountNumber": "",
					},
				),
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
			},
		},
		{
			name: "validate params BankAccountName error",
			args: args{
				p: LinkBankAccountParams{
					UserID:            1,
					BankCode:          "VIB",
					BankAccountNumber: "bank_account_number",
				},
			},
			want: want{
				ac: account.Account{},
				err: exceptions.NewInvalidArgumentError(
					"BankAccountName",
					"bank account name must not be empty",
					map[string]interface{}{
						"BankAccountName": "",
					},
				),
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
			},
		},
		{
			name: "validate params SourceOfFundCode invalid error",
			args: args{
				p: LinkBankAccountParams{
					UserID:            1,
					BankCode:          "VIB",
					BankAccountNumber: "bank_account_number",
					BankAccountName:   "bank_account_name",
				},
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
			},
			want: want{
				ac: account.Account{},
				err: exceptions.NewInvalidArgumentError(
					"SourceOfFundCode",
					"invalid source of fund code",
					map[string]interface{}{
						"BankCode": "bank_code",
					},
				),
			},
		},
		{
			name: "accountRepository.GetAccountsByUserID error",
			args: args{
				p: LinkBankAccountParams{
					UserID:            1,
					BankCode:          "VIB",
					BankAccountNumber: "bank_account_number",
					BankAccountName:   "bank_account_name",
				},
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().GetAccountsByUserID(gomock.Any(), gomock.Any()).Return(nil, uterr)
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().AcquireCreateAccountLockByUserID(gomock.Any(), int64(1)).Return(func() {}, nil)
					return m
				},
			},
			want: want{
				ac:  account.Account{},
				err: uterr,
			},
		},
		{
			name: "Create account success",
			args: args{
				p: LinkBankAccountParams{
					UserID:            1,
					BankCode:          "VIB",
					BankAccountNumber: "bank_account_number",
					BankAccountName:   "bank_account_name",
				},
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().GetAccountsByUserID(gomock.Any(), int64(1)).Return(nil, nil)
					m.EXPECT().CreateAccount(gomock.Any(), account.Account{
						ID:     0,
						UserID: 1,
						Status: account.StatusNormal,
						SourceOfFundData: account.SourceOfFundData{
							IsSourceOfFundItr: implement_bank_account.VibAccount{
								BankAccount: bank_account.BankAccount{
									AccountNumber: "bank_account_number",
									AccountName:   "bank_account_name",
								},
								Status: implement_bank_account.VIBAccountStatusActive,
							},
						},
					}).Return(account.Account{
						ID:     11,
						UserID: 1,
						Status: account.StatusNormal,
						SourceOfFundData: account.SourceOfFundData{
							IsSourceOfFundItr: implement_bank_account.VibAccount{
								BankAccount: bank_account.BankAccount{
									AccountNumber: "bank_account_number",
									AccountName:   "bank_account_name",
								},
								Status: implement_bank_account.VIBAccountStatusActive,
							},
						},
					}, nil)
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().AcquireCreateAccountLockByUserID(gomock.Any(), int64(1)).Return(func() {}, nil)
					return m
				},
			},
			want: want{
				ac: account.Account{
					ID:     11,
					UserID: 1,
					Status: account.StatusNormal,
					SourceOfFundData: account.SourceOfFundData{
						IsSourceOfFundItr: implement_bank_account.VibAccount{
							BankAccount: bank_account.BankAccount{
								AccountNumber: "bank_account_number",
								AccountName:   "bank_account_name",
							},
							Status: implement_bank_account.VIBAccountStatusActive,
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "Create account error with existed account",
			args: args{
				p: LinkBankAccountParams{
					UserID:            1,
					BankCode:          "VIB",
					BankAccountNumber: "bank_account_number",
					BankAccountName:   "bank_account_name",
				},
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().GetAccountsByUserID(gomock.Any(), int64(1)).Return([]account.Account{
						{
							ID:     11,
							UserID: 1,
							Status: account.StatusNormal,
							SourceOfFundData: account.SourceOfFundData{
								IsSourceOfFundItr: implement_bank_account.VibAccount{
									BankAccount: bank_account.BankAccount{
										AccountNumber: "bank_account_number",
										AccountName:   "bank_account_name",
									},
									Status: implement_bank_account.VIBAccountStatusActive,
								},
							},
						},
					}, nil)
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().AcquireCreateAccountLockByUserID(gomock.Any(), int64(1)).Return(func() {}, nil)
					return m
				},
			},
			want: want{
				ac: account.Account{},
				err: exceptions.NewPreconditionError(
					exceptions.PreconditionReasonAccountIsLinked,
					exceptions.SubjectAccount,
					"account is already linked",
					map[string]interface{}{
						"UserID":    int64(1),
						"AccountID": int64(11),
					},
				),
			},
		},
		{
			name: "Create account success with existed unlinked account",
			args: args{
				p: LinkBankAccountParams{
					UserID:            1,
					BankCode:          "VIB",
					BankAccountNumber: "bank_account_number",
					BankAccountName:   "bank_account_name",
				},
			},
			depmock: depmock{
				accountRepository: func(tt *testing.T) *MockaccountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().GetAccountsByUserID(gomock.Any(), int64(1)).Return([]account.Account{
						{
							ID:     11,
							UserID: 1,
							Status: account.StatusUnlinked,
							SourceOfFundData: account.SourceOfFundData{
								IsSourceOfFundItr: implement_bank_account.VibAccount{
									BankAccount: bank_account.BankAccount{
										AccountNumber: "bank_account_number",
										AccountName:   "bank_account_name",
									},
									Status: implement_bank_account.VIBAccountStatusActive,
								},
							},
						},
					}, nil)
					m.EXPECT().CreateAccount(gomock.Any(), account.Account{
						ID:     0,
						UserID: 1,
						Status: account.StatusNormal,
						SourceOfFundData: account.SourceOfFundData{
							IsSourceOfFundItr: implement_bank_account.VibAccount{
								BankAccount: bank_account.BankAccount{
									AccountNumber: "bank_account_number",
									AccountName:   "bank_account_name",
								},
								Status: implement_bank_account.VIBAccountStatusActive,
							},
						},
					}).Return(account.Account{
						ID:     11,
						UserID: 1,
						Status: account.StatusNormal,
						SourceOfFundData: account.SourceOfFundData{
							IsSourceOfFundItr: implement_bank_account.VibAccount{
								BankAccount: bank_account.BankAccount{
									AccountNumber: "bank_account_number",
									AccountName:   "bank_account_name",
								},
								Status: implement_bank_account.VIBAccountStatusActive,
							},
						},
					}, nil)
					return m
				},
				distributeLock: func(tt *testing.T) *MockdistributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().AcquireCreateAccountLockByUserID(gomock.Any(), int64(1)).Return(func() {}, nil)
					return m
				},
			},
			want: want{
				ac: account.Account{
					ID:     11,
					UserID: 1,
					Status: account.StatusNormal,
					SourceOfFundData: account.SourceOfFundData{
						IsSourceOfFundItr: implement_bank_account.VibAccount{
							BankAccount: bank_account.BankAccount{
								AccountNumber: "bank_account_number",
								AccountName:   "bank_account_name",
							},
							Status: implement_bank_account.VIBAccountStatusActive,
						},
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LinkBankAccount{
				accountRepository: tt.depmock.accountRepository(t),
				distributeLock:    tt.depmock.distributeLock(t),
			}
			ac, err := l.Handle(tt.args.ctx, tt.args.p)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.ac, ac)
		})
	}
}
