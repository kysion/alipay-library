package alipay_merchant_v1

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/kysion/alipay-test/alipay_model"
)

// 商家应用相关服务

//// AlipayAuthUserInfo 用户登录信息
//func (c *cAliPay) AlipayAuthUserInfo(ctx context.Context, _ *v1.AlipayAuthUserInfoReq) (StringRes, error) {
//	result, err := alipay_service.AliPay().UserInfoAuth(ctx)
//
//	return StringRes(result), err
//}
//

// H5TradeReq  创建订单，执行支付  https://alipay.kuaimk.com/alipay/h5Pay/6327966408572997/h5Trade
type H5TradeReq struct {
    g.Meta `path:"/:appId/h5Trade" method:"get" summary:"H5支付创建" tags:"阿里云支付"`
    alipay_model.TradeWapPay
}

/*
https://alipay.kuaimk.com/alipay/h5Pay/2021003179623086/h5Trade?totalAmount=0.01&return_url=https%3A%2F%2Fwww.baidu.com&product_code=132313131&subject=测试支付
*/
/*
https://alipay.kuaimk.com/alipay/2021003179623086/gateway.notify"
订单ID：
    6328467214041157
    6328461497139269
    6328484264149061

平台订单ID：
    6328484264149061
*/
