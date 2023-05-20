package alipay_controller

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/api_v1"
	alipay_merchant_v1 "github.com/kysion/alipay-library/alipay_api/alipay_v1/alipay_merchant_v1"
	"github.com/kysion/alipay-library/alipay_service"
)

var MerchantNotify = cMerchantNotify{}

type cMerchantNotify struct{}

// NotifyServices 异步通知地址
func (c *cMerchantNotify) NotifyServices(ctx context.Context, req *alipay_merchant_v1.NotifyServicesReq) (api_v1.StringRes, error) {
	_, err := alipay_service.MerchantNotify().MerchantNotifyServices(ctx)

	return "success", err
}
