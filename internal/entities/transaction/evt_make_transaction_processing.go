package transaction

const EventMakeTransactionProcessingName = EventName("MakeTransactionProcessing")

type MakeTransactionProcessing struct {
	FromStatus Status `json:"from_status"`
	ToStatus   Status `json:"to_status"`
}

func (e MakeTransactionProcessing) EventName() EventName {
	return EventMakeTransactionProcessingName
}
