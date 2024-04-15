package deposit

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

func init() {
	transaction.RegisterTransactionDataEncoder(transaction.TypeDeposit, func(data transaction.IsTransactionDataItr) ([]byte, error) {
		return json.Marshal(data)
	})

	transaction.RegisterTransactionDataDecoder(transaction.TypeDeposit, func(data []byte) (transaction.IsTransactionDataItr, error) {
		var d Deposit
		if err := json.Unmarshal(data, &d); err != nil {
			return nil, err
		}
		return d, nil
	})
}

type partnerJsonData struct {
	SOFType account.SourceOfFundType `json:"sof_type"`
	SOFCode account.SourceOfFundCode `json:"sof_code"`

	Data json.RawMessage `json:"data"`
}

func (s Data) MarshalJSON() ([]byte, error) {
	if s.PartnerData == nil {
		return []byte("{}"), nil
	}
	encoder, err := GetDepositDataEncoder(s.PartnerData.GetSOFType(), s.PartnerData.GetSOFCode())
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit data encoder: %w", err)
	}
	implPartnerData, err := encoder(s.PartnerData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode deposit data: %w", err)
	}
	val := partnerJsonData{
		SOFType: s.PartnerData.GetSOFType(),
		SOFCode: s.PartnerData.GetSOFCode(),
		Data:    implPartnerData,
	}
	bs, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deposit data json data: %w", err)
	}
	return bs, nil
}

func (s *Data) UnmarshalJSON(data []byte) error {
	jsData := partnerJsonData{}
	if err := json.Unmarshal(data, &jsData); err != nil {
		return err
	}
	decoder, err := GetDepositDataDecoder(jsData.SOFType, jsData.SOFCode)
	if err != nil {
		return fmt.Errorf("failed to get deposit data decoder: %w", err)
	}
	implPartnerData, err := decoder(jsData.Data)
	if err != nil {
		return fmt.Errorf("failed to decode deposit data: %w", err)
	}
	*s = Data{
		PartnerData: implPartnerData,
	}

	return nil
}
