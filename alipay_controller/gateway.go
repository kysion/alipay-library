package alipay_controller

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/api_v1"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-test/alipay_service"
	v1 "github.com/kysion/alipay-test/api/alipay_v1"
)

// Gateway 网关
var Gateway = cGateway{}

type cGateway struct{}

type StringRes string

// AliPayServices 商家授权应用，等消息推送，消息通知，通过这个消息  针对B端
func (c *cGateway) AliPayServices(ctx context.Context, req *v1.AliPayServicesReq) (api_v1.BoolRes, error) {
	// fmt.Println("=====Request Begin===========")
	// fmt.Println(g.RequestFromCtx(ctx).GetFormMap())
	// fmt.Println("=====Request End===========")
	result, err := alipay_service.Gateway().GatewayServices(ctx)
	return result != "", err
}

// AliPayCallback C端业务小消息   消费者支付.....
func (c *cGateway) AliPayCallback(ctx context.Context, req *v1.AliPayCallbackReq) (api_v1.BoolRes, error) {

	result, err := alipay_service.Gateway().GatewayCallback(ctx)

	return result != "", err
}

// GetAlipayUserInfo 获取支付宝会员信息，相当于静默登录
func (c *cGateway) GetAlipayUserInfo(ctx context.Context, _ *v1.GetAlipayUserInfoReq) (StringRes, error) {
	// 网页移动应用可以获取用户信息：   回调地址需要换成kuaimk的  https://alipay.kuaimk.com/alipay/1975251903f95826/gateway.callback
	g.RequestFromCtx(ctx).Response.RedirectTo("https://openauth.alipay.com/oauth2/publicAppAuthorize.htm?app_id=2021003179632101&scope=auth_user&redirect_uri=https%3A%2F%2Falipay.kuaimk.com%2Falipay%2F1975251903f95826%2Fgateway.callback")
	// 三方应用和小程序不能获取用户信息，原因如下：
	//     - 这个是小程序应用没有授权回调地址，再换个网页或者、移动应用
	//     - 三方应用是在三方调用情况下使用的，您现在是自调用模式调用 merchant.user.info.share（支付宝会员授权信息查询接口）会因为不是三方调用报此用户不允许自调用

	return "", nil
}
