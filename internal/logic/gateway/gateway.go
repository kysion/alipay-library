package gateway

import (
	"context"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	"github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/alipay_utility"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
)

var (
	// 私钥
	privateData = []byte{}
	// 公钥
	publicCertData = []byte{}
	// 应用公钥
	appCertPublicKeyData = []byte{}
	// 阿里云根证书
	alipayRootCertData = []byte{}
	// 阿里证书公钥
	alipayCertPublicKeyData = []byte{}
)

/*
	阿里网关服务
*/

type sGateway struct {
	CallbackMsgHook base_hook.BaseHook[enum.CallbackMsgType, hook.ServiceMsgHookFunc]

	ServiceNotifyTypeHook base_hook.BaseHook[enum.ServiceNotifyType, hook.ServiceNotifyHookFunc]
}

// GetCallbackMsgHook 返回回调消息Hook对象
func (s *sGateway) GetCallbackMsgHook() *base_hook.BaseHook[enum.CallbackMsgType, hook.ServiceMsgHookFunc] {
	return &s.CallbackMsgHook
}

func (s *sGateway) GetServiceNotifyTypeHook() *base_hook.BaseHook[enum.ServiceNotifyType, hook.ServiceNotifyHookFunc] {
	return &s.ServiceNotifyTypeHook
}

func NewGateway() alipay_service.IGateway {

	return &sGateway{}
}

// GatewayServices 接收消息通知  B端消息
func (s *sGateway) GatewayServices(ctx context.Context) (string, error) {
	urlAppId := g.RequestFromCtx(ctx).Get("appId").String()
	var pathAppId int64
	if urlAppId != "" {
		// 解析AppId
		//pathAppId, _ = strconv.ParseInt(urlAppId, 32, 0)
		pathAppId = alipay_utility.AlipayAppIdDecode(urlAppId)

		if pathAppId == 0 {
			g.RequestFromCtx(ctx).Response.Write("")
			return "参数错误！", nil
		}
	}

	//client, _ := aliyun.NewClient(ctx, gconv.String(pathAppId))
	merchantConfig, _ := alipay_service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(pathAppId))

	var client *aliyun.AliPay
	if merchantConfig.ThirdAppId != "" {
		client, _ = aliyun.NewClient(ctx, gconv.String(pathAppId))
	} else {
		client, _ = aliyun.NewMerchantClient(ctx, gconv.String(pathAppId))
	}

	bm, _ := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)
	fmt.Println(bm)
	msgType := bm.Get("msg_method")
	fmt.Println(msgType)

	fmt.Println("Alipay应用通知类型----------------  ", bm.Get("service"))

	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay应用通知类型： ------- "+bm.Get("service"), "sGateway")

	// 验证应用网关我们直接处理
	if bm.Get("service") == enum.Info.ServiceType.ServiceCheck.Code() {
		s.checkGateway(ctx, client, bm)
	}

	// 通过Hook解决不同的回调类型
	s.ServiceNotifyTypeHook.Iterator(func(key enum.ServiceNotifyType, value hook.ServiceNotifyHookFunc) {
		//if key.Code() == gconv.String(bm.Get("service")) {
		if key.Code() == gconv.String(bm.Get("msg_method")) {
			g.Try(ctx, func(ctx context.Context) {
				sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay应用通知广播： ------- "+bm.Get("source"), "sGateway")
				value(ctx, bm)
			})
		}
	})

	return "", nil
}

// 验证应用网关
func (s *sGateway) checkGateway(ctx context.Context, client *aliyun.AliPay, info gopay.BodyMap) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay验证应用网关： ------- ", "sGateway")

	sign, err := client.GetRsaSign(gopay.BodyMap{
		"success": "true",
	}, "RSA2", "", "xml")

	if err != nil {
		return
	}

	data := gmap.New()

	data.Set("alipay", map[string]interface{}{
		"response": g.Map{
			"success": "true",
		},
		"app_cert_sn": client.AppCertSN,
		"sign":        sign,
		"sign_type":   "RSA2",
	})

	encode, err := gxml.Encode(data.MapStrAny())

	ret := g.RequestFromCtx(ctx).Response

	fmt.Println(string(encode))
	ret.Write("<?xml version=\"1.0\" encoding=\"GBK\"?>")
	ret.Write(string(encode))

	return
}

// GatewayCallback 接收消息回调  C端消息
func (s *sGateway) GatewayCallback(ctx context.Context) (string, error) {

	// 商家的话，先授权，然后获取应用token，存起来

	// 用户的话，直接登录，然后通过code获得token，然后存起来

	request := g.RequestFromCtx(ctx).Request
	fmt.Println(request)

	// 授权之前输入商家信息name -->  签名 -->  --> 签名后存储商家部分数据， --> 自定义授权URL,包含sys_user_id --> 授权，成功的话，根据data找出商家初始数据，然后更新app_auth_token --> 添加第三方平台和用户记录
	bm, err := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)

	fmt.Println("消息回调类型----------------  ", bm.Get("source"))

	data := gopay.BodyMap{
		"grant_type":  "authorization_code",
		"app_id":      bm.Get("app_id"),
		"sys_user_id": bm.Get("sys_user_id"),
		"merchant_id": bm.Get("merchant_id"),
	}

	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay消息回调类型： ------- "+bm.Get("source"), "sGateway")

	// 判断回调的源目标source  HOOK解决switch
	switch bm.Get("source") {
	// 商家应用授权
	case "alipay_app_auth": // 应用授权
		data.Set("code", bm.Get("app_auth_code")) // 商家授权code
	case "alipay_wallet": // 获取用户信息
		data.Set("code", bm.Get("auth_code")) // 用户授权code
	}

	// 通过Hook解决不同的回调类型
	s.CallbackMsgHook.Iterator(func(key enum.CallbackMsgType, value hook.ServiceMsgHookFunc) {
		if key.Code() == gconv.String(bm.Get("source")) {
			g.Try(ctx, func(ctx context.Context) {
				sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay消息回调类型广播： ------- "+bm.Get("source"), "sGateway")
				value(ctx, data)
			})
		}
	})
	// 注意，支付宝回调函数不允许有返回值，不然默认就失败
	return "", err

}
