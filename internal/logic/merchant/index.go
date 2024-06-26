package merchant

import (
	service "github.com/kysion/alipay-library/alipay_service"
)

func init() {

	service.RegisterAppAuth(NewAppAuth())

	service.RegisterMerchantNotify(NewMerchantNotify())

	service.RegisterPayTrade(NewPayTrade())

	service.RegisterMerchantH5Pay(NewMerchantH5Pay())

	service.RegisterMerchantTinyappPay(NewMerchantTinyappPay())

	//service.RegisterWallet(NewWallet())

	service.RegisterMerchantService(NewMerchantService())

	service.RegisterAppVersion(NewAppVersion())

	service.RegisterCertify(NewCertify())

	service.RegisterUserAuth(NewUserAuth())

}
