package transaction

const EventTransactionFailedName = EventName("TransactionFailed")

type EventTransactionFailed struct {
	Transaction Transaction `json:"transaction"`
	FromStatus  Status      `json:"from_status"`
	ToStatus    Status      `json:"to_status"`
}

func (e EventTransactionFailed) EventName() EventName {
	return EventTransactionFailedName
}
