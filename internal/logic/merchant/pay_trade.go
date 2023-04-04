package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_service"
	"strconv"
)

type sPayTrade struct {
}

//
//func init() {
//	service.RegisterPayTrade(NewPayTrade())
//}

func NewPayTrade() *sPayTrade {

	result := &sPayTrade{}

	return result
}

// PayTradeCreate  1、创建交易订单   （AppId的H5是没有的，需要写死，小程序有的 ）
func (s *sPayTrade) PayTradeCreate(ctx context.Context, info *alipay_model.TradeOrder, userId string, notifyFunc ...hook.NotifyHookFunc) (res string, err error) {
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
		res, err = service.MerchantH5Pay().H5TradeCreate(ctx, info, merchantApp, orderInfo, totalAmount)
	}

	return res, err
}
