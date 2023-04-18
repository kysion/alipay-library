package aliyun

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/gopay"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/gopay/pkg/util"
	"github.com/kysion/gopay/pkg/xhttp"
	"github.com/kysion/gopay/pkg/xlog"
	"github.com/kysion/gopay/pkg/xpem"
	"github.com/kysion/gopay/pkg/xrsa"
	"hash"
	"time"
)

/*
	拓展SDK所不具备的
*/

type AliPay struct {
	*alipay.Client
	ThirdConfig     alipay_model.AlipayThirdAppConfig
	MerchantConfig  *alipay_model.AlipayMerchantAppConfig
	privateKey      *rsa.PrivateKey
	aliPayPublicKey *rsa.PublicKey // 支付宝证书公钥内容 alipayCertPublicKey_RSA2.crt

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
		SetCharset(alipay.UTF8). // 设置字符编码，不设置默认 utf-8
		SetSignType(alipay.RSA2). // 设置签名类型，不设置默认 RSA2
		SetReturnUrl(config.AppGatewayUrl). // 设置回调URL
		SetNotifyUrl(config.AppCallbackUrl). // 设置异步通知URL
		SetAppAuthToken(config.AppAuthToken)

	if client.MerchantConfig != nil {
		aliPayClient.SetAppAuthToken(client.MerchantConfig.AppAuthToken)
	}

	key := xrsa.FormatAlipayPrivateKey(config.PrivateKey)
	priKey, err := xpem.DecodePrivateKey([]byte(key))

	client.privateKey = priKey

	//配置公共参数
	aliPayClient.SetCharset("utf-8").
		SetSignType(alipay.RSA2)

	// 自动同步验签（只支持证书模式）
	// 传入 alipayCertPublicKey_RSA2.crt 支付宝证书公钥内容
	aliPayClient.AutoVerifySign([]byte(config.PublicKeyCert))

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

// PostAliPayAPISelfV2 支付宝接口自行实现方法
// 注意：biz_content 需要自行通过bm.SetBodyMap()设置，不设置则没有此参数
// 示例：请参考 client_test.go 的 TestClient_PostAliPayAPISelfV2() 方法
func (s *AliPay) PostAliPayAPISelfV2(ctx context.Context, bm gopay.BodyMap, method string, aliRsp interface{}) (err error) {
	var (
		bs, bodyBs []byte
	)
	// check if there is biz_content
	bz := bm.GetInterface("biz_content")
	if bzBody, ok := bz.(gopay.BodyMap); ok {
		if bodyBs, err = json.Marshal(bzBody); err != nil {
			return fmt.Errorf("json.Marshal(%v)：%w", bzBody, err)
		}
		bm.Set("biz_content", string(bodyBs))
	}

	if bs, err = s.doAliPaySelf(ctx, bm, method); err != nil {
		return err
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return err
	}
	return nil
}

// 向支付宝发送自定义请求
func (s *AliPay) doAliPaySelf(ctx context.Context, bm gopay.BodyMap, method string) (bs []byte, err error) {
	var (
		url, sign string
	)
	bm.Set("method", method)
	// check public parameter
	s.checkPublicParam(bm)

	if bm.GetString("sign") == "" {
		sign, err = s.getRsaSign(bm, bm.GetString("sign_type"), s.privateKey)
		if err != nil {
			return nil, fmt.Errorf("GetRsaSign Error: %w", err)
		}
		bm.Set("sign", sign)
	}

	if s.DebugSwitch == gopay.DebugOn {
		xlog.Debugf("Alipay_Request: %s", bm.JsonBody())
	}

	httpClient := xhttp.NewClient()
	//httpClient.SetBodySize()

	if s.IsProd {
		url = "https://openapi.alipay.com/gateway.do?charset=utf-8"
	} else {
		url = "https://openapi.alipaydev.com/gateway.do?charset=utf-8"
	}
	res, bs, err := httpClient.Type(xhttp.TypeForm).Post(url).SendString(bm.EncodeURLParams()).EndBytes(ctx)
	if err != nil {
		return nil, err
	}
	if s.DebugSwitch == gopay.DebugOn {
		xlog.Debugf("Alipay_Response: %s%d %s%s", xlog.Red, res.StatusCode, xlog.Reset, string(bs))
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	return bs, nil
}

// 公共参数检查
func (s *AliPay) checkPublicParam(bm gopay.BodyMap) {
	bm.Set("format", "JSON").
		Set("charset", s.Charset).
		Set("sign_type", s.SignType).
		Set("version", "1.0").
		Set("timestamp", time.Now().Format(util.TimeLayout))

	if bm.GetString("app_id") == "" && s.AppId != util.NULL {
		bm.Set("app_id", s.AppId)
	}
	if bm.GetString("app_cert_sn") == "" && s.AppCertSN != util.NULL {
		bm.Set("app_cert_sn", s.AppCertSN)
	}
	if bm.GetString("alipay_root_cert_sn") == "" && s.AliPayRootCertSN != util.NULL {
		bm.Set("alipay_root_cert_sn", s.AliPayRootCertSN)
	}
	if bm.GetString("return_url") == "" && s.ReturnUrl != util.NULL {
		bm.Set("return_url", s.ReturnUrl)
	}
	if bm.GetString("notify_url") == "" && s.NotifyUrl != util.NULL {
		bm.Set("notify_url", s.NotifyUrl)
	}
	if bm.GetString("app_auth_token") == "" && s.AppAuthToken != util.NULL {
		bm.Set("app_auth_token", s.AppAuthToken)
	}
}

func (s *AliPay) getRsaSign(bm gopay.BodyMap, signType string, privateKey *rsa.PrivateKey) (sign string, err error) {
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
	signParams := bm.EncodeAliPaySignParams()
	if s.DebugSwitch == gopay.DebugOn {
		xlog.Debugf("Alipay_Request_SignStr: %s", signParams)
	}
	if _, err = h.Write([]byte(signParams)); err != nil {
		return
	}
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, privateKey, hashs, h.Sum(nil)); err != nil {
		return util.NULL, fmt.Errorf("[%w]: %+v", gopay.SignatureErr, err)
	}
	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}
