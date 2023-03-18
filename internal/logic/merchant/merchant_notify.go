package merchant

import (
	"context"
	"github.com/go-pay/gopay/alipay"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_consts"
	dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/pay-share-library/pay_model"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_service"
	"time"
)

/*
	异步通知地址
*/

type sMerchantNotify struct {
	NotifyHook base_hook.BaseHook[hook.NotifyKey, hook.NotifyHookFunc]
	TradeHook  base_hook.BaseHook[hook.TradeHookKey, hook.TradeHookFunc]
}

func init() {
	service.RegisterMerchantNotify(NewMerchantNotify())
}
func NewMerchantNotify() *sMerchantNotify {
	return &sMerchantNotify{}
}

func (s *sMerchantNotify) InstallHook(hookKey hook.NotifyKey, hookFunc hook.NotifyHookFunc) {
	hookKey.HookCreatedAt = *gtime.Now()

	secondAt := gtime.New(alipay_consts.Global.TradeHookExpireAt * gconv.Int64(time.Second))
	hookKey.HookExpireAt = *gtime.New(hookKey.HookCreatedAt.Second() + secondAt.Second())
	s.NotifyHook.InstallHook(hookKey, hookFunc)
}

func (s *sMerchantNotify) InstallTradeHook(hookKey hook.TradeHookKey, hookFunc hook.TradeHookFunc) {
	s.TradeHook.InstallHook(hookKey, hookFunc)
}

// MerchantNotifyServices 异步通知地址  用于接收支付宝推送给商户的支付/退款成功的消息。
func (s *sMerchantNotify) MerchantNotifyServices(ctx context.Context) (string, error) {
	// 1、验签  （自动了）

	// 2、解密 （自动了）

	// 3、获取参数，判断是否交易成功、交易成功记录订单

	//client, err := aliyun.NewClient(ctx, "")

	err := dao.AlipayConsumerConfig.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		bm, _ := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)

		notifyKey := hook.NotifyKey{}

		// 解析消息参数。。。。
		if bm.GetInterface("passback_params") != nil {
			data, has := bm.GetInterface("passback_params").(gmap.Map)
			if has {
				data.GetVar("order_id").String()
				notifyKey.NotifyType = enum.Notify.NotifyType.New(data.GetVar("notify_type").String())
			}
		}

		s.NotifyHook.Iterator(func(key hook.NotifyKey, value hook.NotifyHookFunc) {
			isClean := false
			if key.NotifyType == notifyKey.NotifyType && key.OrderId == notifyKey.OrderId {
				g.Try(ctx, func(ctx context.Context) {
					isClean = value(ctx, kconv.Struct(bm, gmap.Map{}), key)
				})
			}

			s.NotifyHook.UnInstallHook(key, func(filter hook.NotifyKey, conditionKey hook.NotifyKey) bool {
				if key.HookExpireAt.Before(gtime.Now()) {
					return filter == conditionKey
				}
				return isClean && filter == conditionKey
			})
		})

		// 根据订单id找出订单
		order, err := pay_service.Order().GetOrderById(ctx, gconv.Int64(bm["out_trade_no"]))

		bmJson, _ := gjson.Encode(bm)

		// 1. 将交易元数据存储起来 kmk_order
		info := pay_model.UpdateOrderTradeInfo{
			Id:              gconv.Int64(bm["out_trade_no"]),
			PlatformOrderId: gconv.String(bm["trade_no"]), //支付宝交易凭证号。支付宝交易凭证号。
			ConsumerId:      order.ConsumerId,             // 买家支付宝账号对应的支付宝唯一用户号。 gconv.Int64(bm["buyer_id"])
			TradeSource:     gconv.String(bmJson),
		}
		_, err = pay_service.Order().UpdateOrderTradeSource(ctx, &info)
		if err != nil {
			return err
		}

		orderInfo, err := pay_service.Order().GetOrderById(ctx, gconv.Int64(bm["out_trade_no"]))

		// 2.添加定时任务
		gtimer.SetTimeout(ctx, time.Minute*30, func(ctx context.Context) {
			//  判断交易状态，然后修改对应的状态
			var orderState int
			switch bm["trade_status"] {
			case pay_enum.AlipayTrade.TradeStatus.TRADE_SUCCESS.Code():
				// 成功 --> 订单状态为已支付
				orderState = pay_enum.Order.StateType.HavePaid.Code()

			case bm["trade_status"] == pay_enum.AlipayTrade.TradeStatus.TRADE_CLOSED.Code():
				// 交易超时 --> 订单状态为交易超时
				orderState = pay_enum.Order.StateType.PaymentTimeOut.Code()

			case bm["trade_status"] == pay_enum.AlipayTrade.TradeStatus.TRADE_FINISHED.Code():
				// 交易结束，不可退款 --> 订单状态为已完成
				orderState = pay_enum.Order.StateType.DealClose.Code()
			}

			_, err := pay_service.Order().UpdateOrderState(ctx, gconv.Int64(bm["out_trade_no"]), orderState)
			if err != nil {
				return
			}
		})
		if err != nil {
			panic(err)
		}

		// 3. 添加账单account_bill  商家 消费者的账单  业务层Hook
		if bm["trade_status"] == pay_enum.AlipayTrade.TradeStatus.TRADE_SUCCESS.Code() {
			isClean := false

			// Trade发布者
			s.TradeHook.Iterator(func(key hook.TradeHookKey, value hook.TradeHookFunc) {
				if key.AlipayTradeStatus.Code() == pay_enum.AlipayTrade.TradeStatus.TRADE_SUCCESS.Code() {
					su := value(ctx, orderInfo)
					if su == false {
						return
					}
				}

				s.TradeHook.UnInstallHook(key, func(filter hook.TradeHookKey, conditionKey hook.TradeHookKey) bool {
					// 如果超时了，那么久返回true，代表可以删除
					if key.HookExpireAt.Before(gtime.Now()) {
						// 底层的filter和conditionKey是一样的
						return filter == conditionKey
					}
					// 没超时，但是业务层指定了isCLean为true，那么也执行删除
					return isClean && filter == conditionKey
				})
			})

		}

		return nil
		// 积分转移属于账单,,,
	})

	if err != nil {
		return "success", err
	}

	return "success", nil
}
