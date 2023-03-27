package consumer

import "github.com/kysion/base-library/utility/enum"

// ActionEnum 行为：授权、取消授权
type ActionEnum enum.IEnumCode[int]

type action struct {
	Auth   ActionEnum
	UnAuth ActionEnum
}

var ActionType = action{
	Auth:   enum.New[ActionEnum](1, "授权"),
	UnAuth: enum.New[ActionEnum](2, "取消授权"),
}

func (e action) New(code int, description string) ActionEnum {
	if (code&ActionType.Auth.Code()) == ActionType.Auth.Code() ||
		(code&ActionType.UnAuth.Code()) == ActionType.UnAuth.Code() {
		return enum.New[ActionEnum](code, description)
	}
	panic("consumerAction: error")
}
