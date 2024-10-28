package transactions

type ListTransactionHistoriesRequestBody struct {
	OperationType   string  `json:"operation_type" example:"deposit"`
	Amount          float64 `json:"amount" example:"1000"`
	TargetAccountID *string `json:"target_account_id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
}
