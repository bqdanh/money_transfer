package make_transaction

import "github.com/bqdanh/money_transfer/internal/entities/currency"

type MakeDepositTransactionParams struct {
	//RequestID for detect duplicate request, so request id just need check unique in 7 days, that is trace off for performance
	RequestID    string
	UserID       int64
	AccountID    int64
	Amount       currency.Amount
	Descriptions string
	Source       string //optional
}
