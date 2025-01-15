package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/xlog"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/pay-share-library/pay_model"
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
func NewMerchantTinyappPay() service.IMerchantTinyappPay {

	result := &sMerchantTinyappPay{}

	return result
}

// OrderSend 1、发送订单消息   前端完成了
func (s *sMerchantTinyappPay) OrderSend(ctx context.Context) {

}

// TradeCreate  2、小程序创建交易订单
func (s *sMerchantTinyappPay) TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, merchantApp *alipay_model.AlipayMerchantAppConfig, orderInfo *pay_model.OrderRes, totalAmount float32, userId string) (string, error) {
	//appId, _ := strconv.ParseInt(info.AppId, 32, 0)
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------小程序创建交易订单 ------- ", "sMerchantH5Pay")

	client, err := aliyun.NewClient(ctx, merchantApp.AppId)
	notifyUrl := merchantApp.NotifyUrl
	//配置公共参数
	client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetReturnUrl(info.ReturnUrl).
		SetNotifyUrl(notifyUrl)

	orderId := orderInfo.Id // 提交给支付宝的订单Id就是写我们平台数据库中的订单id

	bm := make(gopay.BodyMap)
	bm.Set("subject", info.ProductName).
		Set("buyer_id", userId).
		Set("out_trade_no", orderId).
		Set("total_amount", totalAmount)

	// 创建订单
	// alipay.alipay_trade.create(统一收单交易创建接口)
	aliRsp, err := client.TradeCreate(ctx, bm)
	if err != nil && aliRsp.Response.ErrorResponse.Msg != "Success" {
		if bizErr, ok := alipay.IsBizError(err); ok {
			xlog.Errorf("%s, %s", bizErr.Code, bizErr.Msg)
			// do something
			return "", err
		}
		xlog.Errorf("%s", err)
		return "", err
	}
	xlog.Debug("aliRsp:", *aliRsp)
	xlog.Debug("aliRsp.TradeNo:", aliRsp.Response.TradeNo)

	//g.RequestFromCtx(ctx).Response.WriteJson(aliRsp.Response.TradeNo)

	// 返回具体的订单信息并存储

	return aliRsp.Response.TradeNo, err
}
