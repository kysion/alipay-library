package alipay_hook

import (
    "context"
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
    pay_enum.AlipayTradeStatus
    TradeNo string `json:"tradeNo" dc:"订单交易号"`
}