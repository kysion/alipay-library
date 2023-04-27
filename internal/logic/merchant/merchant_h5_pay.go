package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/gopay"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/gopay/pkg/xlog"
	"github.com/kysion/pay-share-library/pay_model"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_model/pay_hook"
)

type sMerchantH5Pay struct {
	base_hook.BaseHook[pay_enum.OrderStateType, pay_hook.OrderHookFunc]
}

func NewMerchantH5Pay() *sMerchantH5Pay {

	result := &sMerchantH5Pay{}

	return result
}

// InstallHook 安装Hook的时候，如果状态类型为退款中，需要做响应的退款操作，谨防多模块订阅退款状态，产生重复退款
func (s *sMerchantH5Pay) InstallHook(actionType pay_enum.OrderStateType, hookFunc pay_hook.OrderHookFunc) {
	s.BaseHook.InstallHook(actionType, hookFunc)
}

// TradeCreate H5交易下单
func (s *sMerchantH5Pay) TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, merchantApp *alipay_model.AlipayMerchantAppConfig, orderInfo *pay_model.OrderRes, totalAmount float32, userId string) (string, error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------H5交易下单 ------- ", "sMerchantH5Pay")

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)
	notifyUrl := merchantApp.NotifyUrl
	//配置公共参数
	client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetReturnUrl(info.ReturnUrl).
		SetNotifyUrl(notifyUrl)

	orderId := orderInfo.Id // 提交给支付宝的订单Id就是写我们平台数据库中的订单id

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", info.ProductName)
	bm.Set("buyer_id", userId)
	bm.Set("out_trade_no", orderId)
	bm.Set("total_amount", totalAmount)

	//bm.Set("product_code", info.ProductNumber)
	//bm.Set("passback_params", g.Map{ // 可携带数据，在异步通知的的时候会一起回调回来
	//	"notify_type": enum.Notify.NotifyType.PayCallBack.Code(),
	//	"order_id":    orderId,
	//})
	bm.Set("extend_params", g.Map{ // 业务拓展参数
		"royalty_freeze": true, // 资金冻结标识
	})
	bm.Set("app_auth_token", merchantApp.AppAuthToken)
	// 2.统一下单交易创建
	aliRsp, err := client.TradeCreate(ctx, bm)

	if err != nil && aliRsp.Response.ErrorResponse.Msg != "Success" {
		if bizErr, ok := alipay.IsBizError(err); ok {
			xlog.Errorf("%s, %s", bizErr.Code, bizErr.Msg)
			// do something
			return "", err
		}
		xlog.Errorf("%s", err)
		return "", err
	}
	xlog.Debug("aliRsp:", *aliRsp)
	xlog.Debug("aliRsp.TradeNo:", aliRsp.Response.TradeNo)

	return aliRsp.Response.TradeNo, err

}

// H5 支付，返回支付url
func (s *sMerchantH5Pay) H5TradePay(ctx context.Context, info *alipay_model.TradeOrder, merchantApp *alipay_model.AlipayMerchantAppConfig, orderInfo *pay_model.OrderRes, totalAmount float32) (string, error) {
	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.AppId)
	notifyUrl := merchantApp.NotifyUrl
	//配置公共参数
	client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetReturnUrl(info.ReturnUrl).
		SetNotifyUrl(notifyUrl)

	orderId := orderInfo.Id // 提交给支付宝的订单Id就是写我们平台数据库中的订单id

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", info.ProductName)
	bm.Set("out_trade_no", orderId)
	bm.Set("quit_url", notifyUrl)
	bm.Set("total_amount", totalAmount)
	bm.Set("product_code", info.ProductNumber)
	bm.Set("passback_params", g.Map{ // 可携带数据，在异步通知的的时候会一起回调回来
		"notify_type": enum.Notify.NotifyType.PayCallBack.Code(),
		"order_id":    orderId,
	})
	bm.Set("extend_params", g.Map{ // 业务拓展参数
		"royalty_freeze": true, // 资金冻结标识
	})

	// 2.手机网站支付请求
	payUrl, err := client.TradeWapPay(ctx, bm)
	if err != nil {
		xlog.Error("err:", err)
		return "", err
	}
	xlog.Debug("payUrl:", payUrl)

	// 请求重定向到收银台页面
	//g.RequestFromCtx(ctx).Response.RedirectTo(payUrl)

	// 将url返回给前端
	// g.RequestFromCtx(ctx).Response.WriteJson(payUrl)

	// 查询订单并返回tradeNo给前端
	rsp, err := service.PayTrade().QueryOrderInfo(ctx, gconv.String(orderId), merchantApp)

	return rsp.Response.TradeNo, err

	//return "", err

	//g.RequestFromCtx(ctx).Response.Write(orderId)
}

// 2、异步通知 merchant_notify.go
