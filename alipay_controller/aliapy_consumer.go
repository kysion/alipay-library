package alipay_controller

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/api_v1"
	alipay_consumer_v1 "github.com/kysion/alipay-library/alipay_api/alipay_v1/alipay_consumer"
	service "github.com/kysion/alipay-library/alipay_service"
)

var AlipayConsumerConfig = cAlipayConsumerConfig{}

type cAlipayConsumerConfig struct{}

func (c *cAlipayConsumerConfig) AuditConsumer(ctx context.Context, req *alipay_consumer_v1.CertifyReq) (api_v1.StringRes, error) {
	ret, err := service.UserCertity().AuditConsumer(ctx, &req.CertifyInitReq)
	return (api_v1.StringRes)(ret), err
}
