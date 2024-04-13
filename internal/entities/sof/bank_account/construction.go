package bank_account

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

type SourceOfFundBankAccountConstructor func(a BankAccount) (account.SourceOfFundData, error)

var (
	bankSourceOfFundsString = map[account.SourceOfFundCode]account.SourceOfFundCode{}

	sourceOfFundBankAccountMapConstructor = map[account.SourceOfFundCode]SourceOfFundBankAccountConstructor{}
	rwConstructionMutex                   = sync.RWMutex{}
)

func RegisterSourceOfFundBankAccountConstructor(code account.SourceOfFundCode, constructor SourceOfFundBankAccountConstructor) {
	rwConstructionMutex.Lock()
	defer rwConstructionMutex.Unlock()
	_, ok := sourceOfFundBankAccountMapConstructor[code]
	if ok {
		// panic for avoid duplicate register
		panic("source of fund bank account constructor already registered")
	}
	sourceOfFundBankAccountMapConstructor[code] = constructor
}

func RegisterSourceOfFundBankAccount(code account.SourceOfFundCode) {
	rwConstructionMutex.Lock()
	defer rwConstructionMutex.Unlock()
	v, ok := bankSourceOfFundsString[code]
	if ok {
		panic(fmt.Sprintf("source of fund code %s already registered", v))
	}
	bankSourceOfFundsString[code] = code
}

func CreateSourceOfFundBankAccount(code account.SourceOfFundCode, a BankAccount) (account.SourceOfFundData, error) {
	rwConstructionMutex.RLock()
	defer rwConstructionMutex.RUnlock()
	constructor, ok := sourceOfFundBankAccountMapConstructor[code]
	if !ok {
		return account.SourceOfFundData{}, exceptions.NewInvalidArgumentError("SourceOfFundCode", "invalid source of fund code", map[string]interface{}{"code": code})
	}
	return constructor(a)
}

func FromStringToSourceOfFundCode(strCode string) (account.SourceOfFundCode, error) {
	strCode = strings.ToUpper(strCode)
	code := account.SourceOfFundCode(strCode)
	v, ok := bankSourceOfFundsString[code]
	if !ok {
		return "", exceptions.NewInvalidArgumentError("SourceOfFundCode", "invalid source of fund code", map[string]interface{}{"code": code})
	}
	return v, nil
}
