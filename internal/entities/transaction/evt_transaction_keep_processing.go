package transaction

const EventKeepProcessingName = EventName("TransactionKeepProcessing")

type EventTransactionKeepProcessing struct {
	TransData Data `json:"deposit_result"`
}

func (e EventTransactionKeepProcessing) EventName() EventName {
	return EventKeepProcessingName
}
