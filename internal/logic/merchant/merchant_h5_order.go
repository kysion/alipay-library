package merchant

import (
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_model/pay_hook"
)

type sH5Order struct {
	base_hook.BaseHook[pay_enum.OrderStateType, pay_hook.OrderHookFunc]
}

func NewH5Order() *sH5Order {

	result := &sH5Order{}

	return result
}
