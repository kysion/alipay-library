package alipay_hook

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

// SubAccountHookFunc 分账HookFunc （ctx, 参数是下单参数 ） 使用场景：当支付成功后，Hook传递订单数据，然后在业务层创建账单
type SubAccountHookFunc func(ctx context.Context, info *alipay_model.TradeOrderSettleReq) bool

type SubAccountHookInfo struct {
	Key   SubAccountHookKey
	Value ServiceMsgHookFunc
}

type SubAccountHookKey struct {
	SubAccountNo string `json:"tradeNo" dc:"订单交易号"`

	HookCreatedAt gtime.Time `json:"hook_created_at" dc:"Hook创建时间"`
	HookExpireAt  gtime.Time `json:"hook_expire_at" dc:"Hook有效期"`
	// 分账行为
	alipay_enum.SubAccountAction
}
