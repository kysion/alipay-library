package alipay_v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type AliPayServicesReq struct {
	g.Meta `path:"/:appId/gateway.services" method:"post" summary:"阿里云网关消息接收" tags:"阿里云"`
}

type AliPayCallbackReq struct {
	g.Meta `path:"/:appId/gateway.callback" method:"get"  summary:"阿里云网关回调" tags:"阿里云"`
}

//
//type AlipayAuthUserInfoReq struct {
//	g.Meta `path:"/gateway.auth" method:"get" summary:"获取用户授权" tags:"阿里云"`
//}

// GetAlipayUserInfoReq 获取支付宝会员信息，相当于静默登录
type GetAlipayUserInfoReq struct {
	g.Meta `path:"/:appId/gateway.call" method:"get" summary:"获取支付宝会员信息" tags:"阿里云"`
}
