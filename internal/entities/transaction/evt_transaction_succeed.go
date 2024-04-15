package transaction

const EventTransactionSucceedName = EventName("TransactionSucceed")

type EventTransactionSucceed struct {
	Transaction Transaction `json:"transaction"`
	FromStatus  Status      `json:"from_status"`
	ToStatus    Status      `json:"to_status"`
}

func (e EventTransactionSucceed) EventName() EventName {
	return EventTransactionSucceedName
}
