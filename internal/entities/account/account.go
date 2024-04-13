package account

type Account struct {
	ID               int64
	UserID           int64
	Status           Status
	SourceOfFundData SourceOfFundData
}

type Status string

const (
	StatusNormal   = Status("normal")
	StatusLocked   = Status("locked")
	StatusUnlinked = Status("unlinked")
)

type SourceOfFundType string

const (
	SofTypeBankAccount = SourceOfFundType("bank_account")

	//example for another type of source of fund
	SofTypeEWallet   = SourceOfFundType("ewallet")
	SofTypeBankToken = SourceOfFundType("bank_token")
)

// SourceOfFundCode is the code of source of fund: banks: VIB, Vietcombank, Techcombank, etc, finance institutions: LFVN, Momo, ZaloPay, etc
type SourceOfFundCode string

type SourceOfFundData struct {
	IsSourceOfFundItr
}

type IsSourceOfFundItr interface {
	isSourceOfFund()
	GetSourceOfFundCode() SourceOfFundCode
	GetSourceOfFundType() SourceOfFundType
	IsTheSameSof(other IsSourceOfFundItr) bool
}

type IsSourceOfFundImplementMustImport struct {
}

func (b IsSourceOfFundImplementMustImport) isSourceOfFund() {}
