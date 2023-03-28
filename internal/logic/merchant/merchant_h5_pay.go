package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/gopay"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/gopay/pkg/xlog"
	"strconv"

	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_model/pay_hook"
	"github.com/kysion/pay-share-library/pay_service"
)

type sMerchantH5Pay struct {
	base_hook.BaseHook[pay_enum.OrderStateType, pay_hook.OrderHookFunc]
}

func init() {
	service.RegisterMerchantH5Pay(NewMerchantH5Pay())
}

func NewMerchantH5Pay() *sMerchantH5Pay {

	result := &sMerchantH5Pay{}

	return result
}

// InstallHook 安装Hook的时候，如果状态类型为退款中，需要做响应的退款操作，谨防多模块订阅退款状态，产生重复退款
func (s *sMerchantH5Pay) InstallHook(actionType pay_enum.OrderStateType, hookFunc pay_hook.OrderHookFunc) {
	s.BaseHook.InstallHook(actionType, hookFunc)
}

// H5TradeCreate  1、创建交易订单   （AppId的H5是没有的，需要写死，小程序有的 ）
func (s *sMerchantH5Pay) H5TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, notifyFunc ...hook.NotifyHookFunc) {
	// 100分 / 100 = 1元  10 /100
	totalAmount := gconv.Float32(info.Amount) / 100.0

	// 商家AppId解析，获取商家应用，创建阿里支付客户端
	appId, _ := strconv.ParseInt(g.RequestFromCtx(ctx).Get("appId").String(), 32, 0)
	appIdStr := gconv.String(appId)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appIdStr)
	if err != nil {
		return
	}

	sysUser, err := sys_service.SysUser().GetSysUserById(ctx, merchantApp.SysUserId)
	if err != nil {
		return
	}

	info.Order.TradeSourceType = pay_enum.Order.TradeSourceType.Alipay.Code() // 交易源类型
	info.Order.UnionMainId = merchantApp.UnionMainId
	info.Order.UnionMainType = sysUser.Type

	// 支付前创建交易订单，支付后修改交易订单元数据
	orderInfo, err := pay_service.Order().CreateOrder(ctx, &info.Order) // CreatedOrder不能修改订单id
	if err != nil || orderInfo == nil {
		return
	}

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

	// 如果设置了异步通知地址
	if len(notifyFunc) > 0 {
		// 将异步通知中的APPId拿出来，先订阅，收到支付结果通知时，再广播
		service.MerchantNotify().InstallNotifyHook(hook.NotifyKey{
			NotifyType: enum.Notify.NotifyType.PayCallBack,
			OrderId:    gconv.String(orderId),
		}, notifyFunc[0])
	}

	// 2.手机网站支付请求
	payUrl, err := client.TradeWapPay(ctx, bm)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("payUrl:", payUrl)

	g.RequestFromCtx(ctx).Response.RedirectTo(payUrl)

	return
}

// 2、异步通知 merchant_notify.go

// QueryOrderInfo 查询订单
func (s *sMerchantH5Pay) QueryOrderInfo(ctx context.Context, outTradeNo string, merchantAppId string, thirdAppId string, appAuthToken string) {

	client, _ := aliyun.NewClient(ctx, thirdAppId)

	client.SetAppAuthToken(appAuthToken)

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo)

	//查询订单
	aliRsp, err := client.TradeQuery(ctx, bm)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("订单数据:", *aliRsp)
}
