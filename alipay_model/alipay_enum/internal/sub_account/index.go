package sub_account

type subAccount struct {
	SubAccountBindRes  sub
	Action             action
	OperationType      operationType
	TradeSubAccountRes tradeSubAccount
}

var SubAccount = subAccount{
	SubAccountBindRes: SubAccountBindRes,
	Action:            Action,
	OperationType:     OperationType,

	// 分账
	TradeSubAccountRes: TradeSubAccountRes,
}
