package bank_account

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/bqdanh/money_transfer/internal/entities/account"
)

var (
	ErrorInvalidBankAccountData = fmt.Errorf("invalid bank account data")
)

func init() {
	account.RegisterSourceOfFundEncoder(account.SofTypeBankAccount, encodeBankAccountSourceOfFund)
	account.RegisterSourceOfFundDecoder(account.SofTypeBankAccount, decodeBankAccount)
}

type bankAccountJsonData struct {
	SourceOfFundCode account.SourceOfFundCode `json:"sof_code"`
	SofData          json.RawMessage          `json:"sof_data"`
}

type BankAccountDecoder func(bs []byte) (account.IsSourceOfFundItr, error)
type BankAccountEncoder func(account.IsSourceOfFundItr) ([]byte, error)

var (
	bankAccountDecoderMap = map[account.SourceOfFundCode]BankAccountDecoder{}
	bankAccountEncoderMap = map[account.SourceOfFundCode]BankAccountEncoder{}
	rwEncodeDecodeMutex   = sync.RWMutex{}
)

func RegisterBankAccountDecoder(sofCode account.SourceOfFundCode, decoder BankAccountDecoder) {
	rwEncodeDecodeMutex.Lock()
	defer rwEncodeDecodeMutex.Unlock()
	_, ok := bankAccountDecoderMap[sofCode]
	if ok {
		// panic for avoid duplicate register
		panic("bank account decoder already registered")
	}
	bankAccountDecoderMap[sofCode] = decoder
}

func RegisterBankAccountEncoder(sofCode account.SourceOfFundCode, encoder BankAccountEncoder) {
	rwEncodeDecodeMutex.Lock()
	defer rwEncodeDecodeMutex.Unlock()
	_, ok := bankAccountEncoderMap[sofCode]
	if ok {
		// panic for avoid duplicate register
		panic("bank account encoder already registered")
	}
	bankAccountEncoderMap[sofCode] = encoder
}

func getBankAccountDecoder(sofCode account.SourceOfFundCode) (BankAccountDecoder, bool) {
	rwEncodeDecodeMutex.RLock()
	defer rwEncodeDecodeMutex.RUnlock()
	decoder, ok := bankAccountDecoderMap[sofCode]

	return decoder, ok
}

func getBankAccountEncoder(sofCode account.SourceOfFundCode) (BankAccountEncoder, bool) {
	rwEncodeDecodeMutex.RLock()
	defer rwEncodeDecodeMutex.RUnlock()
	encoder, ok := bankAccountEncoderMap[sofCode]

	return encoder, ok

}

func encodeBankAccountSourceOfFund(ac account.IsSourceOfFundItr) ([]byte, error) {
	encoder, ok := getBankAccountEncoder(ac.GetSourceOfFundCode())
	if !ok {
		return nil, fmt.Errorf("bank account source of fund(%s) not registered: %w", ac.GetSourceOfFundCode(), ErrorInvalidBankAccountData)
	}
	bs, err := encoder(ac)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bank account source of fund(%s): %w", ac.GetSourceOfFundCode(), err)
	}
	data := bankAccountJsonData{
		SourceOfFundCode: ac.GetSourceOfFundCode(),
		SofData:          bs,
	}

	bs2, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json data of bank account source of fund: %w", err)
	}
	return bs2, nil
}

func decodeBankAccount(bs []byte) (account.IsSourceOfFundItr, error) {
	var data bankAccountJsonData
	if err := json.Unmarshal(bs, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bank account source of fund: %w", err)
	}
	decoder, ok := getBankAccountDecoder(data.SourceOfFundCode)
	if !ok {
		return nil, fmt.Errorf("bank account source of fund(%s) not registered: %w", data.SourceOfFundCode, ErrorInvalidBankAccountData)
	}

	sofImpl, err := decoder(data.SofData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bank account source of fund(%s) data: %w", data.SourceOfFundCode, err)
	}
	return sofImpl, nil
}
