package transaction

type EventName string

type EventData interface {
	EventName() EventName
}

type Event struct {
	TransactionID   int64
	TransactionType Type
	Version         int
	Name            EventName
	Data            EventData
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
