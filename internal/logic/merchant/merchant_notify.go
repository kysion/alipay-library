package merchant

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_consts"
	dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/pay-share-library/pay_model"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_service"
	"time"
)

/*
	异步通知地址
*/

type sMerchantNotify struct {
	// 异步通知Hook
	NotifyHook base_hook.BaseHook[hook.NotifyKey, hook.NotifyHookFunc]

	// 交易Hook
	TradeHook base_hook.BaseHook[hook.TradeHookKey, hook.TradeHookFunc]

	// 分账Hook (暂时没用到)
	SubAccountHook base_hook.BaseHook[hook.SubAccountHookKey, hook.SubAccountHookFunc]
}

func init() {
	service.RegisterMerchantNotify(NewMerchantNotify())
}

func NewMerchantNotify() *sMerchantNotify {
	return &sMerchantNotify{}
}

// InstallNotifyHook 订阅异步通知Hook
func (s *sMerchantNotify) InstallNotifyHook(hookKey hook.NotifyKey, hookFunc hook.NotifyHookFunc) {
	hookKey.HookCreatedAt = *gtime.Now()

	secondAt := gtime.New(alipay_consts.Global.TradeHookExpireAt * gconv.Int64(time.Second))
	hookKey.HookExpireAt = *gtime.New(hookKey.HookCreatedAt.Second() + secondAt.Second())

	s.NotifyHook.InstallHook(hookKey, hookFunc)
}

// InstallTradeHook 订阅支付Hook
func (s *sMerchantNotify) InstallTradeHook(hookKey hook.TradeHookKey, hookFunc hook.TradeHookFunc) {
	hookKey.HookCreatedAt = *gtime.Now()

	secondAt := gtime.New(alipay_consts.Global.TradeHookExpireAt * gconv.Int64(time.Second))

	hookKey.HookExpireAt = *gtime.New(hookKey.HookCreatedAt.Second() + secondAt.Second())

	s.TradeHook.InstallHook(hookKey, hookFunc)
}

// InstallSubAccountHook 订阅异步通知Hook
func (s *sMerchantNotify) InstallSubAccountHook(hookKey hook.SubAccountHookKey, hookFunc hook.SubAccountHookFunc) {
	hookKey.HookCreatedAt = *gtime.Now()

	secondAt := gtime.New(alipay_consts.Global.TradeHookExpireAt * gconv.Int64(time.Second))
	hookKey.HookExpireAt = *gtime.New(hookKey.HookCreatedAt.Second() + secondAt.Second())

	s.SubAccountHook.InstallHook(hookKey, hookFunc)
}

// MerchantNotifyServices 异步通知地址  用于接收支付宝推送给商户的支付/退款成功的消息。
func (s *sMerchantNotify) MerchantNotifyServices(ctx context.Context) (string, error) {
	err := dao.AlipayConsumerConfig.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		bm, _ := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)
		notifyKey := hook.NotifyKey{}

		// 解析消息参数。。。。
		if bm.GetInterface("passback_params") != nil {
			data := kconv.Struct(bm.GetInterface("passback_params"), &gmap.Map{})

			data.GetVar("order_id").String()

			notifyKey.NotifyType = enum.Notify.NotifyType.New(data.GetVar("notify_type").String())
		}

		// 广播异步通知Hook
		s.NotifyHook.Iterator(func(key hook.NotifyKey, value hook.NotifyHookFunc) {
			isClean := false
			if key.NotifyType == notifyKey.NotifyType {
				if key.OrderId != "" && key.OrderId != notifyKey.OrderId { // 指定id订阅的情况
					return
				}
				g.Try(ctx, func(ctx context.Context) { // 满足条件，Hook调用
					isClean = value(ctx, kconv.Struct(bm, gmap.Map{}), key)
				})
			}

			if key.OrderId != "" {
				s.NotifyHook.UnInstallHook(key, func(filter hook.NotifyKey, conditionKey hook.NotifyKey) bool {
					if key.HookExpireAt.Before(gtime.Now()) {
						return filter == conditionKey
					}
					return isClean && filter == conditionKey
				})
			}
		})

		bmJson, _ := gjson.Encode(bm)

		{
			// 1. 将交易元数据存储起来 kmk_order
			info := pay_model.UpdateOrderTradeInfo{
				Id:              gconv.Int64(bm["out_trade_no"]),
				PlatformOrderId: gconv.String(bm["trade_no"]), //支付宝交易凭证号。支付宝交易凭证号。
				TradeSource:     gconv.String(bmJson),
			}
			_, err := pay_service.Order().UpdateOrderTradeSource(ctx, &info)
			if err != nil {
				return err
			}
		}

		{
			// 2.判断交易状态，然后修改对应的状态
			var orderState int
			switch bm["trade_status"] {
			case pay_enum.AlipayTrade.TradeStatus.TRADE_SUCCESS.Code():
				// 成功 --> 订单状态为已支付
				orderState = pay_enum.Order.StateType.Paymented.Code()

			case bm["trade_status"] == pay_enum.AlipayTrade.TradeStatus.TRADE_CLOSED.Code():
				// 交易超时 --> 订单状态为交易超时
				orderState = pay_enum.Order.StateType.PaymentTimeOut.Code()

			case bm["trade_status"] == pay_enum.AlipayTrade.TradeStatus.TRADE_FINISHED.Code():
				// 交易结束，不可退款 --> 订单状态为已完成
				orderState = pay_enum.Order.StateType.PaymentComplete.Code()
			}

			_, err := pay_service.Order().UpdateOrderState(ctx, gconv.Int64(bm["out_trade_no"]), orderState)
			if err != nil {
				return err
			}
		}

		orderInfo, err := pay_service.Order().GetOrderById(ctx, gconv.Int64(bm["out_trade_no"]))
		if err != nil {
			return err
		}

		// 3. 添加账单account_bill  商家 消费者的账单  业务层Hook
		if bm["trade_status"] == pay_enum.AlipayTrade.TradeStatus.TRADE_SUCCESS.Code() {
			// 4. 分账交易下单结算  需要支付状态为Success的订单

			// a.查询分账关系
			//relationBatch, _ := service.SubAccount().TradeRelationBatchQuery(ctx, gconv.String(bm["auth_app_id"]), gconv.String(bm["out_trade_no"]))
			//if relationBatch.AlipayTradeRoyaltyRelationBatchqueryResponse.ResultCode == enum.SubAccount.SubAccountBindRes.Fail.Code() {
			//	return nil
			//}

			// b.找到分账支出方账户  可选

			// c.组装分账明细信息 + 分账拓展参数

			// 2.次日分账，添加分账的定时任务

			// e.分账通知会发送到应用网关，然后我们判断分账结果，从而创建分账快照
			// alipay.trade.order.settle.notify(交易分账结果通知)  这是我们自己定义的接口吗，不，是应用网关

			// f.上面这些全部写到了具体业务端的定时任务器中操作，先查询所有未分账的账单，然后进行分账，然后请求未分账标记信息

			isClean := false

			// Trade发布广播
			s.TradeHook.Iterator(func(key hook.TradeHookKey, value hook.TradeHookFunc) {
				if key.AlipayTradeStatus.Code() == pay_enum.AlipayTrade.TradeStatus.TRADE_SUCCESS.Code() {
					//var v interface{} = value
					//{
					//	h, ok  := v.(hook.TradeHookFunc[pay_model.OrderRes])
					//	if ok {
					//		h(ctx, orderInfo)
					//	}
					//}
					//
					//
					value(ctx, orderInfo)
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
	})

	if err != nil {
		return "success", err
	}
	g.RequestFromCtx(ctx).Response.Write("success")

	return "success", nil
}
