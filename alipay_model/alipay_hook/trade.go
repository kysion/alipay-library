package alipay_hook

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/kysion/pay-share-library/pay_model"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
)

// TradeHookFunc 支付HookFunc （ctx, 参数是订单） 使用场景：当支付成功后，Hook传递订单数据，然后在业务层创建账单
type TradeHookFunc func(ctx context.Context, info *pay_model.OrderRes) bool

type TradeHookInfo struct {
	Key   TradeHookKey
	Value ServiceMsgHookFunc
}

type TradeHookKey struct {
	TradeNo string `json:"tradeNo" dc:"订单交易号"`

	HookCreatedAt gtime.Time `json:"hook_created_at" dc:"Hook创建时间"`
	HookExpireAt  gtime.Time `json:"hook_expire_at" dc:"Hook有效期"`
	Count         int        `json:"count" dc:"Hook执行次数"`
	// 交易类型
	pay_enum.AlipayTradeStatus
}
