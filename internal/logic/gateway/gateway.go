package gateway

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	enum "github.com/kysion/alipay-test/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-test/alipay_model/alipay_hook"
	"github.com/kysion/alipay-test/internal/logic/internal/aliyun"
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
	base_hook.BaseHook[enum.InfoType, hook.ServiceMsgHookFunc]
}

func (s *sGateway) InstallHook(infoType enum.InfoType, hookFunc hook.ServiceMsgHookFunc) {
	s.BaseHook.InstallHook(infoType, hookFunc)
}

func (s *sGateway) GetHook() base_hook.BaseHook[enum.InfoType, hook.ServiceMsgHookFunc] {
	return s.BaseHook
}

func NewGateway() *sGateway {

	return &sGateway{}
}

// GatewayServices 接收消息通知  B端消息
func (s *sGateway) GatewayServices(ctx context.Context) (string, error) {

	client, err := aliyun.NewClient(ctx, "")
	bm, err := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)

	aliRsp, err := client.OpenAuthTokenAppInviteCreate(ctx, bm)

	fmt.Println(aliRsp)

	g.RequestFromCtx(ctx).Response.Write(aliRsp)
	return aliRsp.Response.TaskPageUrl, err
}

// GatewayCallback 接收消息回调  C端消息
func (s *sGateway) GatewayCallback(ctx context.Context) (string, error) {
	// 商家的话，先授权，然后获取应用token，存起来

	// 用户的话，直接登录，然后通过code获得token，然后存起来

	// 授权之前输入商家信息name -->  签名 -->  --> 签名后存储商家部分数据， --> 自定义授权URL,包含sys_user_id --> 授权，成功的话，根据data找出商家初始数据，然后更新app_auth_token --> 添加第三方平台和用户记录
	bm, err := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)

	data := gopay.BodyMap{
		"grant_type":  "authorization_code",
		"app_id":      bm.Get("app_id"),
		"sys_user_id": bm.Get("sys_user_id"),
	}

	// 判断回调的源目标source  HOOK解决switch
	switch bm.Get("source") {
	// 商家应用授权
	case "alipay_app_auth": // 应用授权
		data.Set("merchant_name", bm.Get("merchant_name"))

		data.Set("code", bm.Get("app_auth_code")) // 商家授权code
	case "alipay_wallet": // 获取用户信息
		data.Set("code", bm.Get("auth_code")) // 用户授权code
	}

	// 通过Hook解决不同的回调类型
	s.Iterator(func(key enum.InfoType, value hook.ServiceMsgHookFunc) {
		if key.Code() == gconv.String(bm.Get("source")) {
			g.Try(ctx, func(ctx context.Context) {
				value(ctx, data)
			})
		}
	})
	// 注意，支付宝回调函数不允许有返回值，不然默认就失败
	return "", err

}
