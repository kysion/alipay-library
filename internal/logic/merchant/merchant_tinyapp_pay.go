package merchant

import (
    "context"
    "github.com/go-pay/gopay"
    "github.com/go-pay/gopay/alipay"
    "github.com/go-pay/gopay/pkg/util"
    "github.com/go-pay/gopay/pkg/xlog"
    "github.com/kysion/alipay-test/alipay_model"
    "github.com/kysion/alipay-test/internal/logic/internal/aliyun"
)

/*
	商户支付相关
*/

// 1、发送订单消息

// 2、创建交易订单

// 3、唤起收银台支付

// 4、异步通知

type sMerchantTinyappPay struct {
}

func NewMerchantTinyappPay() *sMerchantTinyappPay {

    result := &sMerchantTinyappPay{}

    return result
}

// OrderSend 1、发送订单消息
func (s *sMerchantTinyappPay) OrderSend(ctx context.Context) {

}

// TradeCreate  2、创建交易订单
func (s *sMerchantTinyappPay) TradeCreate(ctx context.Context, info *alipay_model.CreateTrade) (aliRsp *alipay_model.TradeCreateResponse, err error) {
    client, _ := aliyun.NewClient(ctx, "")

    // 请求参数  需要传递

    bm := make(gopay.BodyMap)
    bm.Set("subject", "创建订单").
        Set("buyer_id", "2088802095984694").
        Set("out_trade_no", util.RandomString(32)).
        Set("total_amount", "0.01")

    // 创建订单
    // alipay.alipay_trade.create(统一收单交易创建接口)
    client.TradeCreate(ctx, bm)
    if err != nil {
        if bizErr, ok := alipay.IsBizError(err); ok {
            xlog.Errorf("%s, %s", bizErr.Code, bizErr.Msg)
            // do something
            return
        }
        xlog.Errorf("%s", err)
        return
    }
    xlog.Debug("aliRsp:", *aliRsp)
    xlog.Debug("aliRsp.TradeNo:", aliRsp.Response.TradeNo)

    // 返回具体的订单信息并存储

    return nil, err
}
