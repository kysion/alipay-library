package alipay_controller

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/api_v1"
	v1 "github.com/kysion/alipay-library/alipay_api/alipay_v1"
	"github.com/kysion/alipay-library/alipay_service"
)

// Gateway 网关
var Gateway = cGateway{}

type cGateway struct{}

type StringRes string

// AliPayServices 应用通知：针对B端，商家授权应用，等消息推送，消息通知
func (c *cGateway) AliPayServices(ctx context.Context, req *v1.AliPayServicesReq) (api_v1.BoolRes, error) {
	result, err := alipay_service.Gateway().GatewayServices(ctx)
	return result != "", err
}

// GatewayServices 应用网关设置
func (c *cGateway) GatewayServices(ctx context.Context, req *v1.GatewayServicesReq) (api_v1.StringRes, error) {
	alipay_service.Gateway().GatewayServices(ctx)

	//g.RequestFromCtx(ctx).Response.Write("success")

	return "success", nil
}

// AliPayCallback 回调消息：针对C端业务消息   消费者支付.....
func (c *cGateway) AliPayCallback(ctx context.Context, req *v1.AliPayCallbackReq) (api_v1.BoolRes, error) {
	result, err := alipay_service.Gateway().GatewayCallback(ctx)

	return result != "", err
}
