package aliyun

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/pkg/xpem"
	"github.com/go-pay/gopay/pkg/xrsa"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_service"
	"hash"
)

/*
	拓展SDK所不具备的
*/

type AliPay struct {
	*alipay.Client
	ThirdConfig    alipay_model.AlipayThirdAppConfig
	MerchantConfig *alipay_model.AlipayMerchantAppConfig
}

// NewClient 传入各种证书相关文件路径， 初始化客户端对象  appId是商家应用的AppId或者第三方应用的AppId
func NewClient(ctx context.Context, appId string) (client *AliPay, err error) {
	aliPayClient := &alipay.Client{}

	client = &AliPay{}

	if appId == "" {
		//appId = "2021003179681073"
		return nil, sys_service.SysLogs().ErrorSimple(ctx, nil, "非法操作！", "")
	} else {
		client.MerchantConfig, err = alipay_service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)

		if client.MerchantConfig != nil {
			appId = client.MerchantConfig.ThirdAppId
		}
	}
	config, err := alipay_service.ThirdAppConfig().GetThirdAppConfigByAppId(ctx, appId)
	if err != nil {
		return nil, err
	}
	client.ThirdConfig = *config

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
								SetCharset(alipay.UTF8).             // 设置字符编码，不设置默认 utf-8
								SetSignType(alipay.RSA2).            // 设置签名类型，不设置默认 RSA2
								SetReturnUrl(config.AppGatewayUrl).  // 设置返回URL
								SetNotifyUrl(config.AppCallbackUrl). // 设置异步通知URL
								SetAppAuthToken(config.AppAuthToken)

	if client.MerchantConfig != nil {
		aliPayClient.SetAppAuthToken(client.MerchantConfig.AppAuthToken)
	}

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
	client.Client = aliPayClient

	return client, nil
}

// GetSignData 获取数据签名
//func GetSignData(bs []byte, alipayCertSN string) (signData string, err error) {
//	var (
//		str        = string(bs)
//		indexStart = strings.Index(str, `_response":`)
//		indexEnd   int
//	)
//	indexStart = indexStart + 11
//	bsLen := len(str)
//	// 公钥证书模式
//	indexEnd = strings.Index(str, `,"alipay_cert_sn":`)
//	if indexEnd > indexStart && bsLen > indexStart {
//		signData = str[indexStart:indexEnd]
//		return
//	}
//	return gopay.NULL, fmt.Errorf("[%w], value: %s", gopay.GetSignDataErr, str)
//}

// 获取支付宝参数签名
// bm：签名参数
// signType：签名类型，alipay.RSA 或 alipay.RSA2
// privateKey：应用私钥，支持PKCS1和PKCS8
func (s *AliPay) GetRsaSign(bm gopay.BodyMap, signType string, privateKey string, format string) (sign string, err error) {
	if privateKey == "" {
		privateKey = s.ThirdConfig.PrivateKey
	}
	signParams := ""
	key := xrsa.FormatAlipayPrivateKey(privateKey)
	priKey, err := xpem.DecodePrivateKey([]byte(key))

	var (
		h              hash.Hash
		hashs          crypto.Hash
		encryptedBytes []byte
	)
	switch signType {
	case "RSA":
		h = sha1.New()
		hashs = crypto.SHA1
	case "RSA2":
		h = sha256.New()
		hashs = crypto.SHA256
	default:
		h = sha256.New()
		hashs = crypto.SHA256
	}
	format = gstr.ToLower(format)

	if format == "json" {
		signParams = bm.JsonBody()
	} else if format == "xml" {
		byteData, _ := gxml.Encode(bm)
		signParams = string(byteData)
	} else {
		signParams = bm.EncodeURLParams()
	}

	if _, err = h.Write([]byte(signParams)); err != nil {
		return
	}
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, priKey, hashs, h.Sum(nil)); err != nil {
		return util.NULL, fmt.Errorf("[%w]: %+v", gopay.SignatureErr, err)
	}
	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}
