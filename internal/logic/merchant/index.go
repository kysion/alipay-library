package merchant

import (
	service "github.com/kysion/alipay-test/alipay_service"
)

func init() {

	service.RegisterAppAuth(NewAppAuth())

	service.RegisterMerchantNotify(NewMerchantNotify())

	service.RegisterMerchantH5Pay(NewMerchantH5Pay())
	service.RegisterMerchantTinyappPay(NewMerchantTinyappPay())

	service.RegisterWallet(NewWallet())

	service.RegisterMerchantService(NewMerchantService())
}
