package alipay_trade

import "github.com/kysion/base-library/utility/enum"

type TradeStatusEnum enum.IEnumCode[string]

// 交易状态类型
type tradeStatus struct {
    WAIT_BUYER_PAY TradeStatusEnum
    TRADE_CLOSED   TradeStatusEnum
    TRADE_SUCCESS  TradeStatusEnum
    TRADE_FINISHED TradeStatusEnum
}

var TradeStatus = tradeStatus{
    WAIT_BUYER_PAY: enum.New[TradeStatusEnum]("WAIT_BUYER_PAY", "交易创建，等待买家付款。"),
    TRADE_CLOSED:   enum.New[TradeStatusEnum]("TRADE_CLOSED", "未付款交易超时关闭，或支付完成后全额退款。"),
    TRADE_SUCCESS:  enum.New[TradeStatusEnum]("TRADE_SUCCESS", "交易支付成功。"),
    TRADE_FINISHED: enum.New[TradeStatusEnum]("TRADE_FINISHED", "交易结束，不可退款。"),
}
