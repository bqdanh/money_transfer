package deposit

import (
	"fmt"
	"sync"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

type PartnerDataBuilder func(ac account.Account) IsPartnerDepositData
type PartnerDataEncoder func(data IsPartnerDepositData) ([]byte, error)
type PartnerDataDecoder func(bs []byte) (IsPartnerDepositData, error)

var (
	mapDepositDataBuilder = map[string]PartnerDataBuilder{}
	mapDepositEncoder     = map[string]PartnerDataEncoder{}
	mapDepositDecoder     = map[string]PartnerDataDecoder{}

	rwLock sync.RWMutex
)

func RegisterDepositDataBuilder(sofType account.SourceOfFundType, sofCode account.SourceOfFundCode, builder PartnerDataBuilder) {
	rwLock.Lock()
	defer rwLock.Unlock()
	id := fmt.Sprintf("%s_%s", sofType, sofCode)
	if _, ok := mapDepositDataBuilder[id]; ok {
		panic(fmt.Errorf("duplicate deposit data builder for %s", id))
	}
	mapDepositDataBuilder[id] = builder
}

func RegisterDepositDataEncoder(sofType account.SourceOfFundType, sofCode account.SourceOfFundCode, encoder PartnerDataEncoder) {
	rwLock.Lock()
	defer rwLock.Unlock()
	id := fmt.Sprintf("%s_%s", sofType, sofCode)
	if _, ok := mapDepositEncoder[id]; ok {
		panic(fmt.Errorf("duplicate deposit data encoder for %s", id))
	}
	mapDepositEncoder[id] = encoder
}

func GetDepositDataEncoder(sofType account.SourceOfFundType, sofCode account.SourceOfFundCode) (PartnerDataEncoder, error) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	id := fmt.Sprintf("%s_%s", sofType, sofCode)
	encoder, ok := mapDepositEncoder[id]
	if !ok {
		return nil, exceptions.NewInvalidArgumentError(
			"sof_data",
			fmt.Sprintf("deposit data encoder not found for %s", id),
			map[string]interface{}{"sof_type": sofType, "sof_code": sofCode},
		)
	}
	return encoder, nil
}

func RegisterDepositDataDecoder(sofType account.SourceOfFundType, sofCode account.SourceOfFundCode, decoder PartnerDataDecoder) {
	rwLock.Lock()
	defer rwLock.Unlock()
	id := fmt.Sprintf("%s_%s", sofType, sofCode)
	if _, ok := mapDepositDecoder[id]; ok {
		panic(fmt.Errorf("duplicate deposit data decoder for %s", id))
	}
	mapDepositDecoder[id] = decoder
}

func GetDepositDataDecoder(sofType account.SourceOfFundType, sofCode account.SourceOfFundCode) (PartnerDataDecoder, error) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	id := fmt.Sprintf("%s_%s", sofType, sofCode)
	decoder, ok := mapDepositDecoder[id]
	if !ok {
		return nil, exceptions.NewInvalidArgumentError(
			"sof_data",
			fmt.Sprintf("deposit data decoder not found for %s", id),
			map[string]interface{}{"sof_type": sofType, "sof_code": sofCode},
		)
	}
	return decoder, nil
}

func GetDepositDataBuilder(sofType account.SourceOfFundType, sofCode account.SourceOfFundCode) (PartnerDataBuilder, error) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	id := fmt.Sprintf("%s_%s", sofType, sofCode)
	builder, ok := mapDepositDataBuilder[id]
	if !ok {
		return nil, exceptions.NewInvalidArgumentError(
			"sof_data",
			fmt.Sprintf("deposit data builder not found for %s", id),
			map[string]interface{}{"sof_type": sofType, "sof_code": sofCode},
		)
	}
	return builder, nil
}
