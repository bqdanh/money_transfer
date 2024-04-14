package account

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	ErrInvalidSofType = fmt.Errorf("invalid source of fund type")
)

type SourceOfFundEncoder func(IsSourceOfFundItr) ([]byte, error)
type SourceOfFundDecoder func(bs []byte) (IsSourceOfFundItr, error)

var (
	sourceOfFundEncoderMap = map[SourceOfFundType]SourceOfFundEncoder{}
	sourceOfFundDecoderMap = map[SourceOfFundType]SourceOfFundDecoder{}
	rwMutex                = sync.RWMutex{}
)

func RegisterSourceOfFundEncoder(sofType SourceOfFundType, encoder SourceOfFundEncoder) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	_, ok := sourceOfFundEncoderMap[sofType]
	if ok {
		// panic for avoid duplicate register
		panic("source of fund encoder already registered")
	}
	sourceOfFundEncoderMap[sofType] = encoder
}

func RegisterSourceOfFundDecoder(sofType SourceOfFundType, decoder SourceOfFundDecoder) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	_, ok := sourceOfFundDecoderMap[sofType]
	if ok {
		// panic for avoid duplicate register
		panic("source of fund decoder already registered")
	}
	sourceOfFundDecoderMap[sofType] = decoder
}

func getSourceOfFundEncoder(sofType SourceOfFundType) (SourceOfFundEncoder, bool) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	encoder, ok := sourceOfFundEncoderMap[sofType]
	return encoder, ok
}

func getSourceOfFundDecoder(sofType SourceOfFundType) (SourceOfFundDecoder, bool) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	decoder, ok := sourceOfFundDecoderMap[sofType]

	return decoder, ok
}

type sofJsonData struct {
	SofType SourceOfFundType `json:"sof_type"`
	SofCode SourceOfFundCode `json:"sof_code"`
	Data    json.RawMessage  `json:"data"`
}

func (s SourceOfFundData) MarshalJSON() ([]byte, error) {
	encoder, ok := getSourceOfFundEncoder(s.GetSourceOfFundType())
	if !ok {
		return nil, fmt.Errorf("source of fund type (%s) not registered: %w", s.GetSourceOfFundType(), ErrInvalidSofType)
	}
	sofImplData, err := encoder(s.IsSourceOfFundItr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode source of fund data: %w", err)
	}
	val := sofJsonData{
		SofType: s.GetSourceOfFundType(),
		SofCode: s.GetSourceOfFundCode(),
		Data:    sofImplData,
	}
	bs, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal source of fund data json data: %w", err)
	}
	return bs, nil
}

func (s *SourceOfFundData) UnmarshalJSON(data []byte) error {
	jsData := sofJsonData{}
	if err := json.Unmarshal(data, &jsData); err != nil {
		return err
	}
	decoder, ok := getSourceOfFundDecoder(jsData.SofType)
	if !ok {
		return fmt.Errorf("source of fund type (%s) not registered: %w", jsData.SofType, ErrInvalidSofType)
	}
	sofImplData, err := decoder(jsData.Data)
	if err != nil {
		return fmt.Errorf("failed to decode source of fund data: %w", err)
	}
	*s = SourceOfFundData{
		IsSourceOfFundItr: sofImplData,
	}

	return nil
}
