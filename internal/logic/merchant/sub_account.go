package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/go-pay/gopay"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/base-library/utility/kconv"
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

// TradeRelationBatchQuery  分账关系查询

// TradeRelationBatchQuery 分账关系查询
func (s *sSubAccount) TradeRelationBatchQuery(ctx context.Context, appId string, outRequestNo string) (*alipay_model.TradeRelationBatchQuery, error) {
	// 商家AppId解析，获取商家应用，创建阿里支付客户端
	// appId, _ := strconv.ParseInt(g.RequestFromCtx(ctx).Get("appId").String(), 32, 0)

	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return nil, err
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)

	// bindReq := alipay_model.TradeRelationBindReq{}
	// bindReq.outRequestNo = gconv.String(orderId) // 外部请求号  == 单号

	// 1.绑定分账关系 （商户与分账方）
	res, err := client.TradeRelationBatchQuery(ctx, gopay.BodyMap{
		"out_request_no": outRequestNo,
	})

	if err != nil || res.Response.ResultCode != enum.SubAccount.SubAccountBindRes.Success.Code() {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "根据外部请求号查询分账方案失败", "分账方案")
	}
	var ret alipay_model.TradeRelationBatchQuery
	gconv.Struct(&res.Response, &ret)

	return &ret, nil
}

// TradeOrderSettleQuery 交易分账查询接口  必须传递第一个参数，或者后两个同时传递
func (s *sSubAccount) TradeOrderSettleQuery(ctx context.Context, appId string, settleNo string, outRequestNo string, tradeNo string) (*alipay_model.TradeOrderSettleQueryRes, error) {
	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return nil, err
	}
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)
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
	gconv.Struct(&res.Response, &ret)

	return &ret, err
}

// TradeOrderSettle 分账交易下单
func (s *sSubAccount) TradeOrderSettle(ctx context.Context, appId string, info alipay_model.TradeOrderSettleReq) (*alipay_model.TradeOrderSettleResponse, error) {
	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return nil, err
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)

	// bindReq.outRequestNo = gconv.String(orderId) // 外部请求号  == 单号
	data := kconv.Struct(info, gopay.BodyMap{})

	// 1.绑定分账关系 （商户与分账方）
	res, err := client.TradeOrderSettle(ctx, data)

	if err != nil || res.Response.TradeNo != enum.SubAccount.SubAccountBindRes.Success.Code() {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "分账统一下单交易结算失败", "分账方案")
	}

	var ret alipay_model.TradeOrderSettleResponse
	gconv.Struct(res, &ret)

	return &ret, nil
}