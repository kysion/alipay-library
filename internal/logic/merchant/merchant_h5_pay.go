package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/gopay"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/gopay/pkg/xlog"
	"github.com/kysion/pay-share-library/pay_model"
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
func (s *sMerchantH5Pay) H5TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, userId string, notifyFunc ...hook.NotifyHookFunc) (res string, err error) {
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

	orderId := orderInfo.Id // 提交给支付宝的订单Id就是写我们平台数据库中的订单id

	// 如果设置了异步通知地址
	if len(notifyFunc) > 0 {
		// 将异步通知中的APPId拿出来，先订阅，收到支付结果通知时，再广播
		service.MerchantNotify().InstallNotifyHook(hook.NotifyKey{
			NotifyType: enum.Notify.NotifyType.PayCallBack,
			OrderId:    gconv.String(orderId),
		}, notifyFunc[0])
	}

	// 判断是小程序还是H5
	if merchantApp.AppType == 1 {
		// 小程序
		res, err = service.MerchantTinyappPay().TradeCreate(ctx, info, merchantApp, orderInfo, totalAmount, userId)
	} else if merchantApp.AppType == 2 {
		// H5
		res, err = s.H5(ctx, info, merchantApp, orderInfo, totalAmount)
	}

	return res, err
}

func (s *sMerchantH5Pay) H5(ctx context.Context, info *alipay_model.TradeOrder, merchantApp *alipay_model.AlipayMerchantAppConfig, orderInfo *pay_model.OrderRes, totalAmount float32) (string, error) {
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
	g.RequestFromCtx(ctx).Response.RedirectTo(payUrl)

	// 将url返回给前端
	// g.RequestFromCtx(ctx).Response.WriteJson(payUrl)

	// 查询订单并返回tradeNo给前端
	//    rsp, err := s.QueryOrderInfo(ctx, gconv.String(orderId), merchantApp)

	//return rsp.Response.TradeNo, err

	return "", err

	//g.RequestFromCtx(ctx).Response.Write(orderId)
}

// 2、异步通知 merchant_notify.go

// QueryOrderInfo 查询订单
func (s *sMerchantH5Pay) QueryOrderInfo(ctx context.Context, outTradeNo string, merchantApp *alipay_model.AlipayMerchantAppConfig) (aliRsp *alipay.TradeQueryResponse, err error) {
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)

	client.SetAppAuthToken(merchantApp.AppAuthToken)

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo)

	//查询订单
	aliRsp, err = client.TradeQuery(ctx, bm)
	if err != nil && aliRsp.Response.ErrorResponse.Msg != "Success" {
		xlog.Error("err:", err)
		return
	}

	//g.RequestFromCtx(ctx).Response.Write(aliRsp.Response.TradeNo)

	xlog.Debug("订单数据:", *aliRsp)

	return aliRsp, err
}
