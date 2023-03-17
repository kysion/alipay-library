package gateway

import "github.com/kysion/alipay-library/alipay_service"

func init() {
    alipay_service.RegisterGateway(NewGateway())
    alipay_service.RegisterThirdAppConfig(NewThirdAppConfig())
    alipay_service.RegisterMerchantAppConfig(NewMerchantAppConfig())
    alipay_service.RegisterConsumer(NewConsumerConfig())

}
