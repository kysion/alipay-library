package alipay_controller

import (
	"context"
	alipay_consumer_v1 "github.com/kysion/alipay-library/alipay_api/alipay_v1/alipay_consumer"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
)

var AlipayConsumerConfig = cAlipayConsumerConfig{}

type cAlipayConsumerConfig struct{}

func (c *cAlipayConsumerConfig) AuditConsumer(ctx context.Context, req *alipay_consumer_v1.CertifyReq) (*alipay_model.UserCertifyOpenRes, error) {
	ret, err := service.UserCertity().AuditConsumer(ctx, &req.CertifyInitReq)
	return ret, err
}

func (c *cAlipayConsumerConfig) AuditConsumerResponse(ctx context.Context, req *alipay_consumer_v1.CertifyResponseReq) (*alipay_model.UserCertifyOpenQueryRes, error) {
	ret, err := service.UserCertity().AuditConsumerResponse(ctx, req.CertifyId)
	return ret, err
}
