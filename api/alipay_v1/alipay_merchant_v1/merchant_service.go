package alipay_merchant_v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-test/alipay_model"
)

// 商家应用相关服务

// H5TradeReq  创建订单，执行支付  https://alipay.kuaimk.com/alipay/h5Pay/6327966408572997/h5Trade
type H5TradeReq struct {
	g.Meta `path:"/:appId/h5Trade" method:"get" summary:"H5支付创建" tags:"阿里云支付"`
	alipay_model.TradeOrder
}

/*
https://alipay.kuaimk.com/alipay/2021003179623086/gateway.notify"
订单ID：
    6328467214041157
    6328461497139269
    6328484264149061

平台订单ID：
    6328484264149061
*/
