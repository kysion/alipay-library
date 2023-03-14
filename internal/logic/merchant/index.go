package merchant

import (
    service "github.com/kysion/alipay-test/alipay_service"
    "github.com/kysion/alipay-test/internal/logic/merchant/app_auth"
    "github.com/kysion/alipay-test/internal/logic/merchant/h5_pay"
    "github.com/kysion/alipay-test/internal/logic/merchant/tinyapp_pay"
    "github.com/kysion/alipay-test/internal/logic/merchant/wallet"
)

func init() {

    service.RegisterAppAuth(app_auth.NewAppAuth())

    service.RegisterMerchantNotify(h5_pay.NewMerchantNotify())

    service.RegisterMerchantH5Pay(h5_pay.NewMerchantH5Pay())
    service.RegisterMerchantPay(tinyapp_pay.NewMerchantTinyappPay())
    service.RegisterWallet(wallet.NewWallet())

    service.RegisterMerchantService(NewMerchantService())
}
