package alipay_merchant_v1

import "github.com/gogf/gf/v2/frame/g"

// NotifyServicesReq 商家端的异步通知
type NotifyServicesReq struct {
	g.Meta `path:"/:appId/gateway.notify" method:"post" summary:"异步通知" tags:"Alipay"`
}
