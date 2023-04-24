package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/gopay"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_model/pay_hook"
)

type sSubAccount struct {
	base_hook.BaseHook[pay_enum.OrderStateType, pay_hook.OrderHookFunc]
}

func init() {
	service.RegisterSubAccount(NewSubAccount())
}

func NewSubAccount() *sSubAccount {

	result := &sSubAccount{}

	return result
}

// TradeRelationBind 分账关系绑定
func (s *sSubAccount) TradeRelationBind(ctx context.Context, appId int64, info *alipay_model.TradeRelationBindReq) (bool, error) {
	// 商家AppId解析，获取商家应用，创建阿里支付客户端
	// appId, _ := strconv.ParseInt(g.RequestFromCtx(ctx).Get("appId").String(), 32, 0)

	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay分账关系绑定 ------- ", "sSubAccount")

	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return false, err
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)

	// bindReq := alipay_model.TradeRelationBindReq{}
	// bindReq.OutRequestNo = gconv.String(orderId) // 外部请求号  == 单号

	data := kconv.Struct(info, &gopay.BodyMap{})
	data.Set("app_auth_token", merchantApp.AppAuthToken)

	// 1.绑定分账关系 （商户与分账方）
	bindRes, _ := client.TradeRelationBind(ctx, *data)
	if bindRes.Response.ResultCode != enum.SubAccount.SubAccountBindRes.Success.Code() {
		return false, err
	}

	return true, nil
}

// TradeRelationUnbind  分账关系解绑
func (s *sSubAccount) TradeRelationUnbind(ctx context.Context, appId string, info *alipay_model.TradeRelationBindReq) (*alipay_model.TradeRelationUnbindResponse, error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay分账关系解绑 ------- ", "sSubAccount")

	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return nil, err
	}
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)
	client.SetAppAuthToken(merchantApp.AppAuthToken)

	res, err := client.TradeRelationUnbind(ctx, gopay.BodyMap{
		"receiver_list":  info.ReceiverList,
		"out_request_no": info.OutRequestNo,
	})
	var ret alipay_model.TradeRelationUnbindResponse
	gconv.Struct(res, &ret)

	return &ret, err
}

// TradeRelationBatchQuery 分账关系查询
func (s *sSubAccount) TradeRelationBatchQuery(ctx context.Context, appId string, outRequestNo string) (*alipay_model.TradeRoyaltyRelationQueryRes, error) {
	// 商家AppId解析，获取商家应用，创建阿里支付客户端
	//appId, _ := strconv.ParseInt(g.RequestFromCtx(ctx).Get("appId").String(), 32, 0)

	appIdStr := gconv.String(appId) // 这里是第三方的APPID
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return nil, err
	}

	// 通过商家中的第三方应用的AppId创建客户端,，但是认证Token需要是商家的
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)
	client.SetAppAuthToken(merchantApp.AppAuthToken)

	// bindReq := alipay_model.TradeRelationBindReq{}
	// bindReq.outRequestNo = gconv.String(orderId) // 外部请求号  == 单号

	// 1.查询绑定分账关系 （商户与分账方）
	res, _ := client.TradeRelationBatchQuery(ctx, gopay.BodyMap{
		"out_request_no": outRequestNo,
		"app_auth_token": merchantApp.AppAuthToken,
	})
	// 6391922844237893    6391922844237893  out_request_no -> 6391981705789509

	if res.AlipayTradeRoyaltyRelationBatchqueryResponse.ResultCode != enum.SubAccount.SubAccountBindRes.Success.Code() {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "根据外部请求号查询分账方案失败", "分账方案")
	}
	ret := alipay_model.TradeRoyaltyRelationQueryRes{}
	gconv.Struct(res, &ret)

	return &ret, nil
}

// TradeOrderSettleQuery 交易分账查询接口  必须传递第一个参数，或者后两个同时传递
func (s *sSubAccount) TradeOrderSettleQuery(ctx context.Context, appId string, settleNo string, outRequestNo string, tradeNo string) (*alipay_model.TradeOrderSettleQueryRes, error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay交易分账查询接口 ------- ", "sSubAccount")

	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return nil, err
	}
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)
	client.SetAppAuthToken(merchantApp.AppAuthToken)

	/*
	   settle_no	    特殊可选	   支付宝分账请求单号，传入该字段，无需再传外部请求号和支付宝交易号
	   out_request_no	特殊可选	   外部请求号，需要和支付宝交易号一起传入
	   trade_no 	    特殊可选     外部请求号，需要和支付宝交易号一起传入
	*/
	res, err := client.TradeOrderSettleQuery(ctx, gopay.BodyMap{
		"settle_no":      settleNo,
		"out_request_no": outRequestNo,
		"trade_no":       tradeNo,
	})
	var ret alipay_model.TradeOrderSettleQueryRes
	gconv.Struct(res.Response, &ret)

	return &ret, err
}

// TradeOrderSettle 分账交易下单
func (s *sSubAccount) TradeOrderSettle(ctx context.Context, appId string, info alipay_model.TradeOrderSettleReq) (*alipay_model.TradeOrderSettleResponse, error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay分账交易下单 ------- ", "sSubAccount")

	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return nil, err
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)
	client.SetAppAuthToken(merchantApp.AppAuthToken)

	// bindReq.outRequestNo = gconv.String(orderId) // 外部请求号  == 单号
	data := kconv.Struct(info, &gopay.BodyMap{})
	data.Set("app_auth_token", merchantApp.AppAuthToken)

	// 1.分账下单 （商户与分账方）
	res, err := client.TradeOrderSettle(ctx, *data)

	if res.Response.ErrorResponse.Code != enum.SubAccount.TradeSubAccountRes.Success.Code() {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "分账统一下单交易结算失败", "分账方案")
	}

	var ret alipay_model.TradeOrderSettleResponse
	gconv.Struct(res, &ret)

	return &ret, nil
}
