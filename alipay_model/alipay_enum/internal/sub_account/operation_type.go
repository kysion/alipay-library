package sub_account

import "github.com/kysion/base-library/utility/enum"

// OperationTypeEnum  分账操作类型：replenish(补差)、replenish_refund(退补差)、transfer(分账)、transfer_refund(退分账)
type OperationTypeEnum enum.IEnumCode[string]

type operationType struct {
	Replenish       OperationTypeEnum
	ReplenishRefund OperationTypeEnum
	Transfer        OperationTypeEnum
	TransferRefund  OperationTypeEnum
}

var OperationType = operationType{
	Replenish:       enum.New[OperationTypeEnum]("replenish", "补差"),
	ReplenishRefund: enum.New[OperationTypeEnum]("replenish_refund", "退补差"),
	Transfer:        enum.New[OperationTypeEnum]("transfer", "分账"),
	TransferRefund:  enum.New[OperationTypeEnum]("transfer_refund", "退分账"),
}

func (e operationType) New(code string, description string) OperationTypeEnum {
	if code == e.Replenish.Code() {
		return e.Replenish
	}
	if code == e.ReplenishRefund.Code() {
		return e.ReplenishRefund
	}
	if code == e.Transfer.Code() {
		return e.Transfer
	}
	if code == e.TransferRefund.Code() {
		return e.TransferRefund
	}
	return enum.New[OperationTypeEnum](code, description)
}

// replenish(补差)、replenish_refund(退补差)、transfer(分账)、transfer_refund(退分账)
