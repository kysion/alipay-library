package alipay_model

// 订单相关  -- 但是我们自己有创建我们的订单表  kmk_order

type TradeAppPay struct {
    Subject     string  `json:"subject" dc:"订单标题"`
    OutTradeNo  string  `json:"outTradeNo" c:"商户网站唯一订单号"`
    TotalAmount float32 `json:"totalAmount"  c:"订单总金额，单位为元"`
}

type TradeWapPay struct {
    ReturnUrl string `json:"return_url" dc:"交易结束后的返回地址"`
    Subject   string `json:"subject" dc:"订单标题"`
    // OutTradeNo  string  `json:"outTradeNo" dc:"商户网站唯一订单号"`
    TotalAmount float32 `json:"totalAmount" dc:"实际交易金额" v:"required#交易金额不能为空"`
    ProductCode string  `json:"productCode" dc:"商品编号"`
}

// NotifyRequest 异步通知返回参数
type NotifyRequest struct {
    TradeNo string `json:"trade_no,omitempty"`

    TradeStatus string `json:"trade_status,omitempty"`

    TotalAmount string `json:"total_amount,omitempty"`

    PassbackParams string `json:"passback_params" dc:"公共回传参数，如果请求时传递了该参数，则返回给商家时会在异步通知时将该参数原样返回。"`
}