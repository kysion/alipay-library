package sub_account

type subAccount struct {
	SubAccountBindRes sub
	OperationType     operationType
}

var SubAccount = subAccount{
	SubAccountBindRes: SubAccountBindRes,
	OperationType:     OperationType,

	// 分账
	TradeSubAccountRes: TradeSubAccountRes,
}
