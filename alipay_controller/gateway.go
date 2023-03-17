package alipay_controller

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/api_v1"
	"github.com/kysion/alipay-library/alipay_service"
	v1 "github.com/kysion/alipay-library/api/alipay_v1"
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
