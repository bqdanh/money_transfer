package transaction

type EventName string

type EventData interface {
	EventName() EventName
}

type Event struct {
	TransactionID   int64     `json:"transaction_id"`
	TransactionType Type      `json:"transaction_type"`
	Version         int       `json:"version"`
	Name            EventName `json:"name"`
	Data            EventData `json:"data"`
}

func NewEvent(transactionID int64, transactionType Type, version int, data EventData) Event {
	return Event{
		TransactionID:   transactionID,
		TransactionType: transactionType,
		Version:         version,
		Name:            data.EventName(),
		Data:            data,
	}
}
