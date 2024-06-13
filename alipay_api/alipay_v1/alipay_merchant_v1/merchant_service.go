package alipay_merchant_v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model"
)

// 商家应用相关服务

// H5TradeReq  创建订单，执行支付  https://alipay.kuaimk.com/alipay/pay/h5Pay/6327966408572997/h5Trade
type H5TradeReq struct {
	g.Meta `path:"/:appId/h5Trade" method:"get" summary:"H5支付创建" tags:"Alipay商户服务"`
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

// AuthMerchantAppReq 商户授权,内部调起授权页面
type AuthMerchantAppReq struct {
	g.Meta    ` path:"/:appId/gateway.auth"  method:"get" summary:"商户授权" tags:"Alipay商户服务"`
	SysUserId string
}

// GetAlipayUserInfoReq 获取支付宝会员信息，相当于静默登录
type GetAlipayUserInfoReq struct {
	g.Meta `path:"/:appId/gateway.call" method:"get" summary:"获取支付宝会员信息" tags:"Alipay商户服务"`
}

// GetUserInfoByAuthCodeReq 根据认证码获取会员信息
type GetUserInfoByAuthCodeReq struct {
	g.Meta `path:"/getUserInfoByAuthCode" method:"get" summary:"根据认证code获取用户信息" tags:"Alipay商户服务"`

	AuthCode string `json:"auth_code" dc:"第三方平台用户唯一标识ID"`

	AppId string `json:"appId" dc:"商户应用id"`
}
