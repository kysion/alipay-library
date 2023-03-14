package alipay_model

// CreateTrade 统一下单请求参数
type CreateTrade struct {
	Subject     string  `json:"subject" dc:"订单标题"`
	BuyerId     string  `json:"buyer_id" dc:"买家支付宝用户ID。"`
	OutTradeNo  string  `json:"out_trade_no" dc:"商户订单号。"`
	TotalAmount float32 `json:"total_amount" dc:"订单总金额。"`
	ProductCode string  `json:"product_code" dc:"商家和支付宝签约的产品码。"`
}

// ===================================================

type TradeCreateResponse struct {
	Response     *TradeCreate `json:"alipay_trade_create_response"`
	AlipayCertSn string       `json:"alipay_cert_sn,omitempty"`
	SignData     string       `json:"-"`
	Sign         string       `json:"sign" dc:"签名"`
}

type TradeCreate struct {
	ErrorResponse
	TradeNo    string `json:"trade_no,omitempty" dc:"商户订单号"`
	OutTradeNo string `json:"out_trade_no,omitempty" dc:"支付宝交易号"`
}

type ErrorResponse struct {
	Code    string `json:"code" dc:"网关返回码"`
	Msg     string `json:"msg"  dc:"网关返回码描述"`
	SubCode string `json:"sub_code,omitempty" dc:"业务返回码"`
	SubMsg  string `json:"sub_msg,omitempty" dc:"业务返回码描述"`
}
