package alipay_model

// TradeRelationBindReq 分账绑定请求参数Req
type TradeRelationBindReq struct { // List
	OutRequestNo string         `json:"out_request_no" dc:"外部请求号" v:"required#外部请求号不能为空"` // 外部请求号 == 订单ID
	ReceiverList []ReceiverList `json:"receiver_list" dc:"分账关系绑定参数"`
}

type ReceiverList struct {
	Type          string `json:"type" dc:"分账接收方方类型。"`                   // 硬编码：loginName  可以是手机号或者邮箱
	Account       string `json:"account" dc:"分账接收方账号"`                  // account.Account_number
	Name          string `json:"name" dc:"分账接收方真实姓名。"`                  // accuont.name
	Memo          string `json:"memo" dc:"分账关系描述"`                      // 硬编码：充电佣金收入
	LoginName     string `json:"login_name" dc:"当前userId对应的支付宝登录号"`     // 邮箱或者phone手机号
	BindLoginName string `json:"bind_login_name" dc:"分账收款方绑定时的支付宝登录号。"` // 同上 可选

	//Amount float32 `json:"amount" dc:"分账金额"` // 分账金额
}

/*
   分账绑定请求响应参数 result_code:
       SUCCESS：分账关系绑定成功；
       FAIL：分账关系绑定失败。
*/

// ====================分账关系绑定Res===============================
type TradeRelationBindResponse struct {
	Response     *TradeRelationBind `json:"alipay_trade_royalty_relation_bind_response"`
	AlipayCertSn string             `json:"alipay_cert_sn,omitempty"`
	SignData     string             `json:"-"`
	Sign         string             `json:"sign"`
}

type TradeRelationBind struct {
	ErrorResponse
	ResultCode string `json:"result_code" dc:"状态码：SUCCESS和FAIL"`
}

// ====================分账关系查询Res===============================
//type TradeRelationBatchQueryResponse struct {
//	Response     *TradeRelationBatchQuery `json:"alipay_trade_order_settle_response"`
//	AlipayCertSn string                   `json:"alipay_cert_sn,omitempty" dc:"证书"`
//	SignData     string                   `json:"-" dc:"签名"`
//	Sign         string                   `json:"sign" dc:"签名"`
//}
//
//type TradeRelationBatchQuery struct {
//	ErrorResponse
//	ResultCode      string          `json:"result_code" dc:"状态码：SUCCESS和FAIL"`
//	ReceiverList    []*ReceiverList `json:"receiver_list" dc:"分账收款方列表详情"`
//	TotalPageNum    int             `json:"total_page_num" dc:"总页数"`
//	TotalRecordNum  int             `json:"total_record_num" dc:"分账关系记录总数"`
//	CurrentPageNum  int             `json:"current_page_num" dc:"当前页数"`
//	CurrentPageSize int             `json:"current_page_size" dc:"当前页面大小"`
//}

// TradeRoyaltyRelationQueryRes 查询分账关系返回值
type TradeRoyaltyRelationQueryRes struct {
	AlipayTradeRoyaltyRelationBatchqueryResponse AlipayTradeRoyaltyRelationBatchqueryResponse `json:"alipay_trade_royalty_relation_batchquery_response"`
	AlipayCertSn                                 string                                       `json:"alipay_cert_sn"`
	Sign                                         string                                       `json:"sign"`
}

//	type ReceiverList struct {
//		Account       string `json:"account"`
//		BindLoginName string `json:"bind_login_name"`
//		LoginName     string `json:"login_name"`
//		Memo          string `json:"memo"`
//		Type          string `json:"type"`
//	}
type AlipayTradeRoyaltyRelationBatchqueryResponse struct {
	Code            string         `json:"code"`
	Msg             string         `json:"msg"`
	CurrentPageNum  int            `json:"current_page_num"`
	CurrentPageSize int            `json:"current_page_size"`
	ReceiverList    []ReceiverList `json:"receiver_list"`
	ResultCode      string         `json:"result_code"`
	TotalPageNum    int            `json:"total_page_num"`
	TotalRecordNum  int            `json:"total_record_num"`
}

// ==================分账关系解绑Res=================================
type TradeRelationUnbindResponse struct {
	Response     *TradeRelationBind `json:"alipay_trade_royalty_relation_unbind_response"`
	AlipayCertSn string             `json:"alipay_cert_sn,omitempty"`
	SignData     string             `json:"-"`
	Sign         string             `json:"sign"`
}

// ====================分账交易下单Req===============================
type TradeOrderSettleReq struct {
	OutRequestNo      string               `json:"out_request_no" dc:"结算请求流水号，由商家自定义"`
	TradeNo           string               `json:"trade_no" dc:"支付宝订单号"`
	RoyaltyParameters []RoyaltyParameters  `json:"royalty_parameters" dc:"分账明细信息。"`
	OperatorId        string               `json:"operator_id" dc:"操作员 ID，商家自定义操作员编号。"`
	ExtendParams      []SettleExtendParams `json:"extend_params" dc:"分账结算业务扩展参数"`
	RoyaltyMode       string               `json:"royalty_mode" dc:"分账模式，目前有两种分账同步执行sync，分账异步执行async，不传默认同步执"`
}

type RoyaltyParameters struct {
	RoyaltyType  string  `json:"royalty_type" dc:"分账类型."`
	TransOut     string  `json:"trans_out" dc:"支出方账户。"`
	TransOutType string  `json:"trans_out_type" dc:"支出方账户类型。"`
	TransInType  string  `json:"trans_in_type" dc:"收入方账户类型。"`
	TransIn      string  `json:"trans_in" dc:"收入方账户。"`
	Amount       float32 `json:"amount" dc:"分账的金额，单位为元"`
	Desc         string  `json:"desc" dc:"分账描述"`
	RoyaltyScene string  `json:"royalty_scene" dc:"可选值：达人佣金、平台服务费、技术服务费、其他"`
	TransInName  string  `json:"trans_in_name" dc:"分账收款方姓名，"`
}

type SettleExtendParams struct {
	RoyaltyFinish string `json:"royalty_finish" dc:"代表该交易分账是否完结，可选值：true/false，默认值为false。true：代表分账完结，则本次分账处理完成后会把该笔交易的剩余冻结金额全额解冻。false：代表分账未完结"`
}

// ====================分账交易下单Res===============================
type TradeOrderSettleResponse struct {
	Response     *TradeOrderSettle `json:"alipay_trade_order_settle_response"`
	AlipayCertSn string            `json:"alipay_cert_sn,omitempty"`
	SignData     string            `json:"-"`
	Sign         string            `json:"sign"`
}

type TradeOrderSettle struct {
	ErrorResponse
	TradeNo  string `json:"trade_no,omitempty" dc:"支付宝交易编号"`
	SettleNo string `json:"settle_no" dc:"支付宝分账单号，可以根据该单号查询单次分账请求执行结果"`
}

// ===================分账交易查询Req================================
type TradeOrderSettleQueryReq struct {
	SettleNo     string `json:"settle_no" dc:"支付宝分账请求单号，传入了这个则无需传递下面两个参数"`
	OutRequestNo string `json:"out_request_no" dc:"外部请求号，需要和支付宝交易号一起传入"`
	TradeNo      string `json:"trade_no" dc:"支付宝交易号，需要和外部请求号一起传入"`
}

// ===================分账交易Res================================
type TradeOrderSettleQueryRes struct {
	Response     *TradeOrderSettleQuery `json:"alipay_trade_order_settle_response"`
	AlipayCertSn string                 `json:"alipay_cert_sn,omitempty" dc:"证书"`
	SignData     string                 `json:"-"`
	Sign         string                 `json:"sign" dc:"签名"`
}

type TradeOrderSettleQuery struct {
	ErrorResponse
	OutTradeNo        string           `json:"out_request_no" dc:"商户分账请求单号"`
	OperationDt       string           `json:"operation_dt" dc:"分账受理时间"`
	RoyaltyDetailList []*RoyaltyDetail `json:"royalty_detail_list" dc:"分账明细"`
}

type RoyaltyDetail struct {
	OperationType string `json:"operation_type" dc:"分账操作类型：replenish(补差)、replenish_refund(退补差)、transfer(分账)、transfer_refund(退分账)"`
	ExecuteDt     string `json:"execute_dt" dc:"分账执行时间"`
	TransOut      string `json:"trans_out" dc:"分账转出账号，只有在operation_type为replenish(补差),transfer_refund(退分账)类型才返回该字段或"`
	TransOutType  string `json:"trans_out_type" dc:"分账转出账号类型"`
	TransIn       string `json:"trans_in" dc:"分账转入账号"`
	TransInType   string `json:"trans_in_type" dc:"分账转入账号类型"`
	Amount        string `json:"amount" dc:"分账金额"`
	State         string `json:"state" dc:"分账状态：SUCCESS成功、FAIL失败"`
	ErrorCode     string `json:"error_code" dc:"分账失败错误码"`
	ErrorDesc     string `json:"error_desc" dc:"分账错误描述信息"`
}
