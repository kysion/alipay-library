package h5_pay

import (
    "context"
    "github.com/go-pay/gopay/alipay"
    "github.com/gogf/gf/v2/container/gmap"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gcron"
    "github.com/gogf/gf/v2/os/gtime"
    "github.com/gogf/gf/v2/util/gconv"
    "github.com/kysion/alipay-test/alipay_consts"
    enum "github.com/kysion/alipay-test/alipay_model/alipay_enum"
    hook "github.com/kysion/alipay-test/alipay_model/alipay_hook"
    "github.com/kysion/base-library/base_hook"
    "github.com/kysion/base-library/utility/kconv"
    "time"
)

/*
	异步通知地址
*/

type sMerchantNotify struct {
    base_hook.BaseHook[hook.NotifyKey, hook.NotifyHookFunc]
}

func NewMerchantNotify() *sMerchantNotify {
    return &sMerchantNotify{}
}

func (s *sMerchantNotify) InstallHook(hookKey hook.NotifyKey, hookFunc hook.NotifyHookFunc) {
    hookKey.HookCreatedAt = *gtime.Now()

    secondAt := gtime.New(alipay_consts.Global.TradeHookExpireAt * gconv.Int64(time.Second))
    hookKey.HookExpireAt = *gtime.New(hookKey.HookCreatedAt.Second() + secondAt.Second())
    s.BaseHook.InstallHook(hookKey, hookFunc)
}

// MerchantNotifyServices 异步通知地址  用于接收支付宝推送给商户的支付/退款成功的消息。
func (s *sMerchantNotify) MerchantNotifyServices(ctx context.Context) (string, error) {
    // 1、验签  （自动了）

    // 2、解密

    // 3、获取参数，判断是否交易成功、交易成功记录订单

    //client, err := aliyun.NewClient(ctx, "")
    bm, _ := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)

    notifyKey := hook.NotifyKey{}

    // 解析消息参数。。。。
    if bm.GetInterface("passback_params") != nil {
        //passback_params := bm.Get("passback_params")
        data, has := bm.GetInterface("passback_params").(gmap.Map)
        if has {
            data.GetVar("order_id").String()
            notifyKey.NotifyType = enum.Notify.NotifyType.New(data.GetVar("notify_type").String())
        }
    }

    s.Iterator(func(key hook.NotifyKey, value hook.NotifyHookFunc) {
        isClean := false
        if key.NotifyType == notifyKey.NotifyType && key.OrderId == notifyKey.OrderId {
            g.Try(ctx, func(ctx context.Context) {
                isClean = value(ctx, kconv.Struct(bm, gmap.Map{}), key)
            })
        }

        s.BaseHook.UnInstallHook(key, func(filter hook.NotifyKey, conditionKey hook.NotifyKey) bool {
            if key.HookExpireAt.Before(gtime.Now()) {
                return filter == conditionKey
            }
            return isClean && filter == conditionKey
        })
    })

    //  需要补充：将交易元数据存储起来 kmk_order

    // 添加定时任务
    _, err := gcron.Add(ctx, "0 30 * * * *", func(ctx context.Context) {

        // 调用共享层的修改状态方法
        g.Log().Print(ctx, "Every hour on the half hour")
    })
    if err != nil {
        panic(err)
    }

    return "", nil
}
