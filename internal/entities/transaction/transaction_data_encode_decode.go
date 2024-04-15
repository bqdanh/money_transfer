package transaction

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	ErrInvalidTransactionType = fmt.Errorf("invalid transaction type")
)

type TransactionDataEncoder func(IsTransactionDataItr) ([]byte, error)
type TransactionDataDecoder func(bs []byte) (IsTransactionDataItr, error)

var (
	TransactionDataEncoderMap = map[Type]TransactionDataEncoder{}
	TransactionDataDecoderMap = map[Type]TransactionDataDecoder{}
	rwMutex                   = sync.RWMutex{}
)

func RegisterTransactionDataEncoder(transType Type, encoder TransactionDataEncoder) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	_, ok := TransactionDataEncoderMap[transType]
	if ok {
		// panic for avoid duplicate register
		panic("transaction encoder already registered")
	}
	TransactionDataEncoderMap[transType] = encoder
}

func RegisterTransactionDataDecoder(transType Type, decoder TransactionDataDecoder) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	_, ok := TransactionDataDecoderMap[transType]
	if ok {
		// panic for avoid duplicate register
		panic("transaction data decoder already registered")
	}
	TransactionDataDecoderMap[transType] = decoder
}

func getTransactionDataEncoder(transType Type) (TransactionDataEncoder, bool) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	encoder, ok := TransactionDataEncoderMap[transType]
	return encoder, ok
}

func getTransactionDataDecoder(transType Type) (TransactionDataDecoder, bool) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	decoder, ok := TransactionDataDecoderMap[transType]

	return decoder, ok
}

type transactionJsonData struct {
	TransType Type            `json:"trans_type"`
	Data      json.RawMessage `json:"data"`
}

func (s Data) MarshalJSON() ([]byte, error) {
	if s.IsTransactionDataItr == nil {
		return []byte("{}"), nil
	}
	encoder, ok := getTransactionDataEncoder(s.GetType())
	if !ok {
		return nil, fmt.Errorf("transaction type (%s) not registered: %w", s.GetType(), ErrInvalidTransactionType)
	}
	implTransData, err := encoder(s.IsTransactionDataItr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode transaction data: %w", err)
	}
	val := transactionJsonData{
		TransType: s.GetType(),
		Data:      implTransData,
	}
	bs, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction data json data: %w", err)
	}
	return bs, nil
}

func (s *Data) UnmarshalJSON(data []byte) error {
	jsData := transactionJsonData{}
	if err := json.Unmarshal(data, &jsData); err != nil {
		return err
	}
	decoder, ok := getTransactionDataDecoder(jsData.TransType)
	if !ok {
		return fmt.Errorf("transaction data (%s) not registered: %w", jsData.TransType, ErrInvalidTransactionType)
	}
	implTransData, err := decoder(jsData.Data)
	if err != nil {
		return fmt.Errorf("failed to decode transaction data: %w", err)
	}
	*s = Data{
		IsTransactionDataItr: implTransData,
	}

	return nil
}
