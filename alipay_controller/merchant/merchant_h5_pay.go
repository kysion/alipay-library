package merchant_controller

import (
	"context"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/api_v1"
	"github.com/gogf/gf/v2/container/gmap"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/api/alipay_v1/alipay_merchant_v1"
)

var MerchantH5Pay = merchantH5Pay{}

type merchantH5Pay struct{}

func (c *merchantH5Pay) H5TradeCreate(ctx context.Context, req *alipay_merchant_v1.H5TradeReq) (api_v1.StringRes, error) {

	// 创建交易订单，生成二维码，根据二维码调起收银台
	service.MerchantH5Pay().H5TradeCreate(ctx, &req.TradeOrder, c.notifyHookFunc)

	return "", nil
}

func (c *merchantH5Pay) notifyHookFunc(ctx context.Context, info gmap.Map, hookInfo hook.NotifyKey) bool {

	// 在此进行订单支付相关判断和信息更新

	fmt.Println("我是自定义回调哦！！！！！！！！！")

	return false
	// 布尔值如果返回true,那么就会忽略和此订单相关的通知

}

// QueryOrderInfo 根据trade_no平台交易ID 或者是orderId订单id查询订单数据
//func (c *merchantService) QueryOrderInfo(ctx context.Context, outTradeNo string, merchantAppId string, thirdAppId string, appAuthToken string) {
//
//}
