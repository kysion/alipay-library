package aliyun

import (
	"context"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_service"
)

/*
	拓展SDK所不具备的
*/

// NewClient 传入各种证书相关文件路径， 初始化客户端对象
func NewClient(ctx context.Context, appId string) (*alipay.Client, error) {
	aliPayClient := &alipay.Client{}

	config := &alipay_model.AlipayThirdAppConfig{}

	if appId == "" {
		appId = "2021003179681073"
	}

	if appId != "" {
		merchantConfig, err := alipay_service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
		if err != nil {
			return nil, err
		}

		appId = merchantConfig.ThirdAppId
	}

	config, err := alipay_service.ThirdAppConfig().GetThirdAppConfigByAppId(ctx, appId)
	if err != nil {
		return nil, err
	}

	//global := alipay_consts.Global
	// 微信：拿到token、每个请求都需要进行携带签名这些

	// 1、初始化支付宝客户端并做配置(appid：应用ID、privateKey：应用私钥，支持PKCS1和PKCS8、isProd：是否是正式环境)
	aliPayClient, err = alipay.NewClient(config.AppId, config.PrivateKey, true)

	if err != nil {
		xlog.Error(err)
		return nil, err
	}

	// 自定义配置http请求接收返回结果body大小，默认 10MB
	// client.SetBodySize() // 没有特殊需求，可忽略此配置

	// 打开Debug开关，输出日志，默认关闭
	aliPayClient.DebugSwitch = gopay.DebugOn

	// 设置支付宝请求 公共参数
	//    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	aliPayClient.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
								SetCharset(alipay.UTF8).            // 设置字符编码，不设置默认 utf-8
								SetSignType(alipay.RSA2).           // 设置签名类型，不设置默认 RSA2
								SetReturnUrl(config.AppGatewayUrl). // 设置返回URL
								SetNotifyUrl(config.AppCallbackUrl) // 设置异步通知URL

	//配置公共参数
	aliPayClient.SetCharset("utf-8").
		SetSignType(alipay.RSA2)

	// 自动同步验签（只支持证书模式）
	// 传入 alipayCertPublicKey_RSA2.crt 内容
	aliPayClient.AutoVerifySign([]byte(config.AppPublicCertKey))

	// 证书路径(应用公钥证书路径、 支付宝根证书文件路径、 支付宝公钥证书文件路径)
	// err = aliPayClient.SetCertSnByPath(config.AppPublicCertKey, config.AlipayRootCertPublicKey, config.PublicKeyCert)

	// 证书内容(应用公钥证书文件内容、支付宝根证书文件内容、支付宝公钥证书文件内容)
	err = aliPayClient.SetCertSnByContent([]byte(config.AppPublicCertKey), []byte(config.AlipayRootCertPublicKey), []byte(config.PublicKeyCert))
	if err != nil {
		xlog.Debug("SetCertSn:", err)
		return nil, err
	}

	return aliPayClient, nil
}
