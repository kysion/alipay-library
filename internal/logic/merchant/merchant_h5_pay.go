package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kuaimk/kmk-share-library/share_model"
	"github.com/kuaimk/kmk-share-library/share_model/share_enum"
	"github.com/kuaimk/kmk-share-library/share_model/share_hook"
	"github.com/kuaimk/kmk-share-library/share_service"
	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
)

type sMerchantH5Pay struct {
	base_hook.BaseHook[share_enum.OrderStateType, share_hook.OrderHookFunc]
}

func init() {
	service.RegisterMerchantH5Pay(NewMerchantH5Pay())
}

func NewMerchantH5Pay() *sMerchantH5Pay {

	result := &sMerchantH5Pay{}

	return result
}

// InstallHook 安装Hook的时候，如果状态类型为退款中，需要做响应的退款操作，谨防多模块订阅退款状态，产生重复退款
func (s *sMerchantH5Pay) InstallHook(actionType share_enum.OrderStateType, hookFunc share_hook.OrderHookFunc) {
	s.BaseHook.InstallHook(actionType, hookFunc)
}

// H5TradeCreate  1、创建交易订单   （AppId的H5是没有的，需要写死，小程序有的 ）
func (s *sMerchantH5Pay) H5TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, notifyFunc ...hook.NotifyHookFunc) {
	// 需要补充：创建我们平台的订单  --> 之后发起支付宝支付请求  --> 支付成功后  --> 异步通知的地方进行修改改订单的交易元数据和第三方订单id trade_no

	user := sys_service.SysSession().Get(ctx).JwtClaimsUser

	// 1.先拿出用户享有的优惠策略，然后校验是否可以在该产品上使用。

	// 2.根据产品编号查询是否具备xx优惠策略， 然后该消费者的优惠策略是否能使用在该产品上面， 有的话然后将对应的优惠填上去，但是现阶段的充电业务是没有啥优惠的。。。

	// 3.查看产品自身优惠策略，如果具备，叠加优惠金额，

	// 4. 计算实际交易金额
	amount := info.Order.OrderAmount

	// 100分 / 100 = 1元  10 /100
	total_amount := gconv.Float32(info.Order.OrderAmount) / 100.0

	// 5.根据UserId查询消费者财务账号的余额是否足够，现阶段没消费者钱包，所以不进行此操作
	// beforeBalance := "用户钱包消费前余额"
	// afterBalance := "用户钱包消费后余额"

	// 6.构建订单信息，填写相应金额

	// 7. 创建订单

	// 8. 支付宝支付...

	// 9. 异步通知记录支付结果及元数据

	// 查询消费者财务账号

	order := share_model.Order{
		ConsumerId:      user.Id,                                        // 消费者ID
		InOutType:       info.InOutType,                                 // 支出
		TradeSourceType: share_enum.Order.TradeSourceType.Alipay.Code(), // 交易源类型
		Amount:          amount,                                         // 实际成交的交易金额，分为单位 1*100 = 100分 0.01*100 = 1分
		CouponAmount:    0,                                              // 优惠金额
		CouponConfig:    "",                                             // 优惠减免金额
		OrderAmount:     info.Order.OrderAmount,                         // 订单金额
		BeforeBalance:   0,                                              // 交易前的余额
		AfterBalance:    0,                                              // 交易后的金额
		ProductName:     info.ProductName,                               // 产品名称
		TradeScene:      info.TradeScene,                                // 交易场景
		ProductNumber:   info.ProductNumber,                             // 产品编号
		UnionMainId:     user.UnionMainId,                               // 关联主体
	}

	// 支付前创建交易订单，支付后修改交易订单元数据
	orderInfo, err := share_service.Order().CreateOrder(ctx, &order)
	if err != nil || orderInfo == nil {
		return
	}

	// 商家AppId
	appId := g.RequestFromCtx(ctx).Get("appId").String()
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil {
		return
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)

	notifyUrl := "https://alipay.kuaimk.com/alipay/" + appId + "/gateway.notify"

	//配置公共参数
	client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetReturnUrl(info.ReturnUrl).
		SetNotifyUrl(notifyUrl).
		SetAppAuthToken(merchantApp.AppAuthToken)

	orderId := orderInfo.Id // 提交给支付宝的订单Id就是写我们平台数据库中的订单id

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", info.ProductName)
	bm.Set("out_trade_no", orderId)
	bm.Set("quit_url", notifyUrl)
	bm.Set("total_amount", total_amount)
	bm.Set("product_code", info.ProductNumber)
	bm.Set("passback_params", g.Map{ // 可携带数据，在哟不通知的的时候会一起回调回来
		"notify_type": enum.Notify.NotifyType.PayCallBack.Code(),
		"order_id":    orderId,
	})

	// 如果设置了异步通知地址
	if len(notifyFunc) > 0 {
		// 将异步通知中的APPId拿出来，
		service.MerchantNotify().InstallHook(hook.NotifyKey{
			NotifyType: enum.Notify.NotifyType.PayCallBack,
			OrderId:    gconv.String(orderId),
		}, notifyFunc[0])
	}

	//手机网站支付请求
	payUrl, err := client.TradeWapPay(ctx, bm)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("payUrl:", payUrl)

	g.RequestFromCtx(ctx).Response.RedirectTo(payUrl)

	return
}

// 2、异步通知 notify

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
