package sub_account

type subAccount struct {
	SubAccountBindRes sub
	Action            action
	OperationType     operationType
}

var SubAccount = subAccount{
	SubAccountBindRes: SubAccountBindRes,
	Action:            Action,
	OperationType:     OperationType,
}
