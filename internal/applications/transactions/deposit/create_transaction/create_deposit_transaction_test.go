package create_transaction

import (
	"context"
	"fmt"
	"testing"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account/implement_bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/currency"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit/bank_account_vib"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var _ = bank_account_vib.DepositData{}

func TestCreateDepositTransaction_Handle(t *testing.T) {
	var errUTForce = fmt.Errorf("ut-force")
	var utAccountNormal = account.Account{
		ID:     1,
		UserID: 11,
		Status: account.StatusNormal,
		SourceOfFundData: account.SourceOfFundData{
			IsSourceOfFundItr: implement_bank_account.VIBAccount{
				BankAccount: bank_account.BankAccount{
					AccountNumber: "AccountNumber_XYZ",
					AccountName:   "AccountName_XYZ",
				},
				Status: implement_bank_account.VIBAccountStatusActive,
			},
		},
	}
	var utAccountUnlinked = account.Account{
		ID:     1,
		UserID: 11,
		Status: account.StatusUnlinked,
		SourceOfFundData: account.SourceOfFundData{
			IsSourceOfFundItr: implement_bank_account.VIBAccount{
				BankAccount: bank_account.BankAccount{
					AccountNumber: "AccountNumber_XYZ",
					AccountName:   "AccountName_XYZ",
				},
				Status: implement_bank_account.VIBAccountStatusActive,
			},
		},
	}
	var utAccountLock = account.Account{
		ID:     1,
		UserID: 11,
		Status: account.StatusLocked,
		SourceOfFundData: account.SourceOfFundData{
			IsSourceOfFundItr: implement_bank_account.VIBAccount{
				BankAccount: bank_account.BankAccount{
					AccountNumber: "AccountNumber_XYZ",
					AccountName:   "AccountName_XYZ",
				},
				Status: implement_bank_account.VIBAccountStatusActive,
			},
		},
	}

	type fields struct {
		cfg                       Config
		distributeLockFunc        func(tt *testing.T) distributeLock
		accountRepositoryFunc     func(tt *testing.T) accountRepository
		transactionRepositoryFunc func(tt *testing.T) transactionRepository
	}
	type args struct {
		ctx context.Context
		p   CreateDepositTransactionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    transaction.Transaction
		wantErr error
	}{
		{
			name: "validate params with invalid user id",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    0,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewInvalidArgumentError(
				"UserID",
				"user must greater than 0",
				map[string]interface{}{
					"user_id": 0,
				},
			),
		},
		{
			name: "validate params with invalid account id",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 0,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewInvalidArgumentError(
				"AccountID",
				"account must greater than 0",
				map[string]interface{}{
					"account_id": 0,
				},
			),
		},
		{
			name: "validate params with amount less than minimum amount",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   0,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewInvalidArgumentError(
				"Amount",
				fmt.Sprintf("amount must greater than %f", 1.0),
				map[string]interface{}{
					"error":          nil,
					"amount":         0,
					"minimum_amount": 1.0,
				},
			),
		},
		{
			name: "validate params with amount greater than maximum amount",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   1_000_000_000 + 1,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewInvalidArgumentError(
				"Amount",
				fmt.Sprintf("amount must less than %f", 1_000_000_000.0),
				map[string]interface{}{
					"error":  nil,
					"amount": 1_000_000_000 + 1,
				},
			),
		},
		{
			name: "validate params with invalid type of currency",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    111,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.USD,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewInvalidArgumentError(
				"Amount",
				fmt.Sprintf("currency %s is not accepted", currency.USD),
				map[string]interface{}{
					"currency":            currency.USD,
					"currencies_accepted": []string{string(currency.VND)},
				},
			),
		},
		{
			name: "fail get account by id",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(account.Account{}, errUTForce)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    111,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want:    transaction.Transaction{},
			wantErr: errUTForce,
		},
		{
			name: "get account by id with invalid user id",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountNormal, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    22,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewPreconditionError(exceptions.PreconditionReasonPermissionDenied, exceptions.SubjectAccount, "use account of another user", map[string]interface{}{
				"account_id": 1,
				"user_id":    22,
			}),
		},
		{
			name: "fail with account status locked",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountLock, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewPreconditionError(
				exceptions.PreconditionReasonAccountStatusNotReadyForDeposit,
				exceptions.SubjectAccount,
				"account is not ready for deposit",
				map[string]interface{}{
					"account_id": 1,
					"status":     account.StatusLocked,
					"expected":   account.StatusNormal,
				},
			),
		},
		{
			name: "fail with account status unlinked",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountUnlinked, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewPreconditionError(
				exceptions.PreconditionReasonAccountStatusNotReadyForDeposit,
				exceptions.SubjectAccount,
				"account is not ready for deposit",
				map[string]interface{}{
					"account_id": 1,
					"status":     account.StatusUnlinked,
					"expected":   account.StatusNormal,
				},
			),
		},
		{
			name: "fail with account status unlinked",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountUnlinked, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewPreconditionError(
				exceptions.PreconditionReasonAccountStatusNotReadyForDeposit,
				exceptions.SubjectAccount,
				"account is not ready for deposit",
				map[string]interface{}{
					"account_id": 1,
					"status":     account.StatusUnlinked,
					"expected":   account.StatusNormal,
				},
			),
		},
		{
			name: "acquire lock for create deposit transaction failed",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().
						AcquireLockForCreateDepositTransaction(
							gomock.AssignableToTypeOf(context.Background()),
							"request-id", DefaultConfig.LockDuration,
						).
						Return(nil, errUTForce)
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountNormal, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want:    transaction.Transaction{},
			wantErr: errUTForce,
		},
		{
			name: "get transaction by request id is duplicated",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().
						AcquireLockForCreateDepositTransaction(
							gomock.AssignableToTypeOf(context.Background()),
							"request-id", DefaultConfig.LockDuration,
						).
						Return(func() {}, nil)
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountNormal, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					m.EXPECT().GetTransactionByRequestID(gomock.AssignableToTypeOf(context.Background()), "request-id").
						Return(transaction.Transaction{
							ID:      1,
							Account: utAccountNormal,
							Amount: currency.Amount{
								Currency: currency.VND,
								Amount:   10_000,
							},
							Description: "Description",
							Status:      transaction.StatusInit,
							Type:        transaction.TypeDeposit,
							Data: transaction.Data{
								IsTransactionDataItr: deposit.Deposit{
									Source:            "source_xyz",
									BankTransactionID: "",
								},
							},
						}, nil)
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{},
			wantErr: exceptions.NewPreconditionError(exceptions.PreconditionReasonTransactionIsAvailable, exceptions.SubjectTransaction, "request id is duplicated", map[string]interface{}{
				"request_id":     "request-id",
				"transaction_id": int64(1),
			}),
		},
		{
			name: "get transaction got error",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().
						AcquireLockForCreateDepositTransaction(
							gomock.AssignableToTypeOf(context.Background()),
							"request-id", DefaultConfig.LockDuration,
						).
						Return(func() {}, nil)
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountNormal, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					m.EXPECT().GetTransactionByRequestID(gomock.AssignableToTypeOf(context.Background()), "request-id").
						Return(transaction.Transaction{}, errUTForce)
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want:    transaction.Transaction{},
			wantErr: errUTForce,
		},
		{
			name: "persist transaction failed",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().
						AcquireLockForCreateDepositTransaction(
							gomock.AssignableToTypeOf(context.Background()),
							"request-id", DefaultConfig.LockDuration,
						).
						Return(func() {}, nil)
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountNormal, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					m.EXPECT().GetTransactionByRequestID(gomock.AssignableToTypeOf(context.Background()), "request-id").
						Return(transaction.Transaction{}, ErrNotFoundTransaction)
					m.EXPECT().
						CreateTransaction(
							gomock.AssignableToTypeOf(context.Background()),
							transaction.Transaction{
								ID:      0,
								Account: utAccountNormal,
								Amount: currency.Amount{
									Currency: currency.VND,
									Amount:   10_000,
								},
								Description: "descriptions",
								Status:      transaction.StatusInit,
								Type:        transaction.TypeDeposit,
								Data: transaction.Data{
									IsTransactionDataItr: deposit.Deposit{
										Source:            "source",
										BankTransactionID: "",
										Data: deposit.Data{
											PartnerData: bank_account_vib.DepositData{
												Status:           bank_account_vib.DepositStatusInit,
												TransactionRefID: "",
											},
										},
									},
								},
							},
						).
						Return(transaction.Transaction{}, errUTForce)
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want:    transaction.Transaction{},
			wantErr: errUTForce,
		},
		{
			name: "create deposit transaction successfully",
			fields: fields{
				cfg: DefaultConfig,
				distributeLockFunc: func(tt *testing.T) distributeLock {
					m := NewMockdistributeLock(gomock.NewController(tt))
					m.EXPECT().
						AcquireLockForCreateDepositTransaction(
							gomock.AssignableToTypeOf(context.Background()),
							"request-id", DefaultConfig.LockDuration,
						).
						Return(func() {}, nil)
					return m
				},
				accountRepositoryFunc: func(tt *testing.T) accountRepository {
					m := NewMockaccountRepository(gomock.NewController(tt))
					m.EXPECT().
						GetAccountsByID(gomock.AssignableToTypeOf(context.Background()), int64(1)).
						Return(utAccountNormal, nil)
					return m
				},
				transactionRepositoryFunc: func(tt *testing.T) transactionRepository {
					m := NewMocktransactionRepository(gomock.NewController(tt))
					m.EXPECT().GetTransactionByRequestID(gomock.AssignableToTypeOf(context.Background()), "request-id").
						Return(transaction.Transaction{}, ErrNotFoundTransaction)
					m.EXPECT().
						CreateTransaction(
							gomock.AssignableToTypeOf(context.Background()),
							transaction.Transaction{
								ID:      0,
								Account: utAccountNormal,
								Amount: currency.Amount{
									Currency: currency.VND,
									Amount:   10_000,
								},
								Description: "descriptions",
								Status:      transaction.StatusInit,
								Type:        transaction.TypeDeposit,
								Data: transaction.Data{
									IsTransactionDataItr: deposit.Deposit{
										Source:            "source",
										BankTransactionID: "",
										Data: deposit.Data{
											PartnerData: bank_account_vib.DepositData{
												Status:           bank_account_vib.DepositStatusInit,
												TransactionRefID: "",
											},
										},
									},
								},
							},
						).
						Return(transaction.Transaction{
							ID:      1,
							Account: utAccountNormal,
							Amount: currency.Amount{
								Currency: currency.VND,
								Amount:   10_000,
							},
							Description: "descriptions",
							Status:      transaction.StatusInit,
							Type:        transaction.TypeDeposit,
							Data: transaction.Data{
								IsTransactionDataItr: deposit.Deposit{
									Source:            "source",
									BankTransactionID: "",
									Data: deposit.Data{
										PartnerData: bank_account_vib.DepositData{
											Status:           bank_account_vib.DepositStatusInit,
											TransactionRefID: "",
										},
									},
								},
							},
						}, nil)
					return m
				},
			},
			args: args{
				ctx: context.Background(),
				p: CreateDepositTransactionParams{
					RequestID: "request-id",
					UserID:    11,
					AccountID: 1,
					Amount: currency.Amount{
						Currency: currency.VND,
						Amount:   10_000,
					},
					Descriptions: "descriptions",
					Source:       "source",
				},
			},
			want: transaction.Transaction{
				ID:      1,
				Account: utAccountNormal,
				Amount: currency.Amount{
					Currency: currency.VND,
					Amount:   10_000,
				},
				Description: "descriptions",
				Status:      transaction.StatusInit,
				Type:        transaction.TypeDeposit,
				Data: transaction.Data{
					IsTransactionDataItr: deposit.Deposit{
						Source:            "source",
						BankTransactionID: "",
						Data: deposit.Data{
							PartnerData: bank_account_vib.DepositData{
								Status:           bank_account_vib.DepositStatusInit,
								TransactionRefID: "",
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCreateDepositTransaction(tt.fields.cfg, tt.fields.distributeLockFunc(t), tt.fields.accountRepositoryFunc(t), tt.fields.transactionRepositoryFunc(t))
			assert.NoError(t, err)
			got, err := c.Handle(tt.args.ctx, tt.args.p)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
