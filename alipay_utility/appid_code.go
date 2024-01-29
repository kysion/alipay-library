package alipay_utility

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
)

// AlipayAppIdEncode 支付宝 - appId的36进制编码
func AlipayAppIdEncode(alipayAppId int64) (appIdEncode string) {
	appIdEncode = strconv.FormatInt(gconv.Int64(alipayAppId), 36)

	return appIdEncode
}

// AlipayAppIdDecode 支付宝 - appId的36进制解码
func AlipayAppIdDecode(appIdEncode string) (alipayAppId int64) {
	// 解码
	alipayAppId, _ = strconv.ParseInt(appIdEncode, 36, 0)

	return alipayAppId
}

// AlipayAppIdDecodeFormCrx 支付宝 - appId的32进制解码 （从请求的上下文获取）
func AlipayAppIdDecodeFormCrx(ctx context.Context) (alipayAppId int64) {
	// 解码
	alipayAppId, _ = strconv.ParseInt(g.RequestFromCtx(ctx).Get("appId").String(), 36, 0)

	return alipayAppId
}

// Alipay-AppId使用示例：
//encode := share_utility.AlipayAppIdEncode(2017102909591815)
//fmt.Println(encode)
//fmt.Println(share_utility.AlipayAppIdDecode(encode))
