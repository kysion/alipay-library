package alipay_model

// ====================单笔资金转账Req===============================
type FundTransUniTransferReq struct {
	OutBizNo       string    `json:"out_biz_no" dc:"商家侧唯一订单号，由商家自定义。对于不同转账请求"`
	TransAmount    float32   `json:"trans_amount" dc:"订单总金额，单位为元，"`
	ProductCode    string    `json:"product_code" dc:"销售产品码。单笔无密转账固定为 TRANS_ACCOUNT_NO_PWD。"`
	BizScene       string    `json:"biz_scene" dc:"业务场景。单笔无密转账固定为 DIRECT_TRANSFER。"`
	OrderTitle     string    `json:"order_title" dc:"转账业务的标题，用于在支付宝用户的账单里显示。"`
	PayeeInfo      PayeeInfo `json:"payee_info" dc:"收款方信息"`
	Remark         string    `json:"remark" dc:"可选：业务备注。"`
	BusinessParams string    `json:"business_params" dc:"可选：转账业务请求的扩展参数，支持传入的扩展参数如下："`
	/*
	   payer_show_name_use_alias：是否展示付款方别名，可选，收款方在支付宝账单中可见。枚举支持：
	        true：展示别名，将展示商家支付宝在商家中心商户信息> 商户基本信息 页面配置的 商户别名。
	        false：不展示别名。默认为 false
	*/
}

// PayeeInfo 收款方信息
type PayeeInfo struct {
	Identity     string `json:"identity" dc:"参与方的标识 ID，"` // 当 identity_type=ALIPAY_USER_ID 时，填写支付宝用户 UID。示例值：2088123412341234。 当 identity_type=ALIPAY_LOGON_ID 时，填写支付宝登录号。示例值：186xxxxxxxx。
	IdentityType string `json:"identity_type" dc:"参与方的标识类型，"`
	Name         string `json:"name" dc:"可选：参与方真实姓名。如果非空，将校验收款支付宝账号姓名一致性。"`
}

// ====================单笔转账Res===============================
type FundTransUniTransferRes struct {
	Response     *TransUniTransferRes `json:"alipay_fund_trans_uni_transfer_response"`
	AlipayCertSn string               `json:"alipay_cert_sn,omitempty"`
	SignData     string               `json:"-"`
	Sign         string               `json:"sign"`
}

type TransUniTransferRes struct {
	ErrorResponse
	OutBizNo       string `json:"out_biz_no,omitempty" dc:"商户订单号"`
	OrderId        string `json:"order_id,omitempty" dc:"支付宝转账订单号"`
	PayFundOrderId string `json:"pay_fund_order_id,omitempty" dc:"支付宝支付资金流水号"`
	Status         string `json:"status,omitempty" dc:"转账单据状态。 SUCCESS（该笔转账交易成功）：成功； FAIL：失败（具体失败原因请参见error_code以及fail_reason返回值）；"`
	TransDate      string `json:"trans_date,omitempty" dc:"订单支付时间，格式为yyyy-MM-dd HH:mm:s"`
}

// ====================资金转账Res===============================
type FundTransToaccountTransferRes struct {
	Response     *TransToaccountTransfer `json:"alipay_fund_trans_toaccount_transfer_response"`
	AlipayCertSn string                  `json:"alipay_cert_sn,omitempty"`
	SignData     string                  `json:"-"`
	Sign         string                  `json:"sign"`
}

type TransToaccountTransfer struct {
	ErrorResponse
	OutBizNo string `json:"out_biz_no,omitempty"`
	OrderId  string `json:"order_id,omitempty"`
	PayDate  string `json:"pay_date,omitempty"`
}
