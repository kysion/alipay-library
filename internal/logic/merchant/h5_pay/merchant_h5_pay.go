package h5_pay

import (
    "context"
    "github.com/go-pay/gopay"
    "github.com/go-pay/gopay/alipay"
    "github.com/go-pay/gopay/pkg/xlog"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/util/gconv"
    "github.com/kysion/alipay-test/alipay_model"
    enum "github.com/kysion/alipay-test/alipay_model/alipay_enum"
    hook "github.com/kysion/alipay-test/alipay_model/alipay_hook"
    service "github.com/kysion/alipay-test/alipay_service"
    "github.com/kysion/alipay-test/internal/logic/internal/aliyun"
    "github.com/yitter/idgenerator-go/idgen"
)

type sMerchantH5Pay struct {
}

func NewMerchantH5Pay() *sMerchantH5Pay {

    result := &sMerchantH5Pay{}

    return result
}

// H5TradeCreate  1、创建交易订单   （AppId的H5是没有的，需要写死，小程序有的 ）
func (s *sMerchantH5Pay) H5TradeCreate(ctx context.Context, info *alipay_model.TradeWapPay, notifyFunc ...hook.NotifyHookFunc) {
    // 需要补充：创建我们平台的订单  --> 之后发起支付宝支付请求  --> 支付成功后  --> 异步通知的地方进行修改改订单的交易元数据和第三方订单id trade_no

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

    /*
        bm := make(gopay.BodyMap)
       bm.Set("subject", "手机网站测试支付").
               Set("out_trade_no", "GZ201909081743431443").
               Set("quit_url", "https://www.fmm.ink").
               Set("total_amount", "100.00").
               Set("product_code", "QUICK_WAP_WAY")
    */

    orderId := idgen.NextId()
    //请求参数
    bm := make(gopay.BodyMap)
    bm.Set("subject", info.Subject)
    bm.Set("out_trade_no", orderId)
    bm.Set("quit_url", notifyUrl)
    bm.Set("total_amount", info.TotalAmount)
    bm.Set("product_code", info.ProductCode)
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
}

// 2、异步通知

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
