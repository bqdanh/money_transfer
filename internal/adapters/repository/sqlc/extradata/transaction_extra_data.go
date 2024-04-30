package extradata

import "encoding/json"

type TransactionExtraData struct {
	TransactionData json.RawMessage `json:"transaction_data"`
}
