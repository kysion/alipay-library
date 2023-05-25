package internal

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_consts"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/alipay_utility/file"
)

func init() {
	InitGlobal()
	// InitFileData()

}

// InitGlobal 初始化公共对象
func InitGlobal() {
	alipay_consts.Global.RSA2 = g.Cfg().MustGet(context.Background(), "service.RSA2").String()
	alipay_consts.Global.PriPath = g.Cfg().MustGet(context.Background(), "service.priPath").String()
	alipay_consts.Global.PublicCrtPath = g.Cfg().MustGet(context.Background(), "service.publicCrtPath").String()
	alipay_consts.Global.AppCertPublicKeyPath = g.Cfg().MustGet(context.Background(), "service.appCertPublicKeyPath").String()
	alipay_consts.Global.AlipayRootCertPath = g.Cfg().MustGet(context.Background(), "service.alipayRootCertPath").String()
	alipay_consts.Global.AlipayCertPublicKeyPath = g.Cfg().MustGet(context.Background(), "service.alipayCertPublicKeyPath").String()
	alipay_consts.Global.AppId = g.Cfg().MustGet(context.Background(), "service.AppId").String()
	alipay_consts.Global.AppCode = g.Cfg().MustGet(context.Background(), "service.appCode").String()
	alipay_consts.Global.AES = g.Cfg().MustGet(context.Background(), "service.AES").String()
	alipay_consts.Global.CallbackUrl = g.Cfg().MustGet(context.Background(), "service.callbackUrl").String()
	alipay_consts.Global.ReturnUrl = g.Cfg().MustGet(context.Background(), "service.returnUrl").String()

	// 交易Hook失效时间
	alipay_consts.Global.TradeHookExpireAt = g.Cfg().MustGet(context.Background(), "service.tradeHookExpireAt").Int64()

}

func InitFileData() {
	// 加载证书文件
	privateData, _ := file.GetFile(alipay_consts.Global.PriPath)
	publicCertData, _ := file.GetFile(alipay_consts.Global.PublicCrtPath)
	appCertPublicKeyData, _ := file.GetFile(alipay_consts.Global.AppCertPublicKeyPath)
	alipayRootCertData, _ := file.GetFile(alipay_consts.Global.AlipayRootCertPath)
	alipayCertPublicKeyData, _ := file.GetFile(alipay_consts.Global.AlipayCertPublicKeyPath)
	info := alipay_model.UpdateThirdKeyCertReq{
		AppId:                   "2021003179681073",
		PrivateKey:              string(privateData),
		PublicKey:               string(publicCertData),
		PublicKeyCert:           string(alipayCertPublicKeyData),
		AppPublicCertKey:        string(appCertPublicKeyData),
		AlipayRootCertPublicKey: string(alipayRootCertData),
	}
	fmt.Println(info)

	_, err := alipay_service.ThirdAppConfig().UpdateThirdKeyCert(context.Background(), &info)
	if err != nil {
		fmt.Println("证书文件存储失败啦~")
	}
}
