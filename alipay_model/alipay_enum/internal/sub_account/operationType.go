package sub_account

import "github.com/kysion/base-library/utility/enum"

type ActionEnum enum.IEnumCode[int]

type action struct {
	Replenish       ActionEnum
	ReplenishRefund ActionEnum
	Transfer        ActionEnum
	TransferRefund  ActionEnum
}

var Action = action{
	Replenish:       enum.New[ActionEnum](1, "补差"),
	ReplenishRefund: enum.New[ActionEnum](2, "退补差"),
	Transfer:        enum.New[ActionEnum](4, "分账"),
	TransferRefund:  enum.New[ActionEnum](8, "退分账"),
}

func (e action) New(code int, description string) ActionEnum {
	if (code&Action.Replenish.Code()) == Action.Replenish.Code() ||
		(code&Action.ReplenishRefund.Code()) == Action.ReplenishRefund.Code() ||
		(code&Action.Transfer.Code()) == Action.Transfer.Code() ||
		(code&Action.TransferRefund.Code()) == Action.TransferRefund.Code() {
		return enum.New[ActionEnum](code, description)
	} else {
		panic("SubAccount.Action.New: error")
	}
}

// replenish(补差)、replenish_refund(退补差)、transfer(分账)、transfer_refund(退分账)
