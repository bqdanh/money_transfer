package transaction

const EventTransactionCreatedName = EventName("TransactionCreated")

type EventTransactionCreated struct {
	Transaction Transaction `json:"transaction"`
}

func (e EventTransactionCreated) EventName() EventName {
	return EventTransactionCreatedName
}
