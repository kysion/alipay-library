package merchant

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/gopay"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/gopay/pkg/util"
	"github.com/kysion/gopay/pkg/xlog"
	"strconv"
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

func init() {
	service.RegisterMerchantTinyappPay(NewMerchantTinyappPay())
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
	appId, _ := strconv.ParseInt(info.AppId, 32, 0)

	client, _ := aliyun.NewClient(ctx, gconv.String(appId))

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
