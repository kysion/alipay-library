package merchant

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-test/alipay_model"
	enum "github.com/kysion/alipay-test/alipay_model/alipay_enum"
	service "github.com/kysion/alipay-test/alipay_service"
	"github.com/kysion/alipay-test/internal/logic/internal/aliyun"
)

/*
	获取支付宝会员信息等
*/

type sWallet struct {
}

func NewWallet() *sWallet {
	// 初始化文件内容

	result := &sWallet{}

	result.injectHook()
	return result
}

func (s *sWallet) injectHook() {
	service.Gateway().InstallHook(enum.Info.Type.AlipayWallet, s.Wallet)
}

// Wallet 具体服务
func (s *sWallet) Wallet(ctx context.Context, info g.Map) bool {
	client, _ := aliyun.NewClient(ctx, "")

	data := gopay.BodyMap{}
	gconv.Struct(info, &data)

	// 根据AppId获取商家相关配置，包括AppAuthToken
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(info["app_id"]))
	if err != nil || merchantApp == nil {
		return false
	}

	// data.Set("code", data.Get("auth_code"))                            // 用户授权code

	// 这个token是动态的，哪个商家需要获取，appId和appAuthToken就传递对应的
	// client.SetAppAuthToken("202303BB6460ae32bcd14f20be7c0c5eff86dE11") // 商家Token
	client.SetAppAuthToken(merchantApp.AppAuthToken)

	// 1.auth_code换Token
	token, _ := client.SystemOauthToken(ctx, data)

	// token获取支付宝会员授权信息查询接口
	aliRsp, _ := client.UserInfoShare(ctx, token.Response.AccessToken)

	fmt.Println(token)
	fmt.Println(aliRsp)

	authInfo := map[string]interface{}{
		"userId":             aliRsp.Response.UserId,
		"avatar":             aliRsp.Response.Avatar,
		"province":           aliRsp.Response.Province,
		"city":               aliRsp.Response.City,
		"nickName":           aliRsp.Response.NickName,
		"isStudentCertified": aliRsp.Response.IsStudentCertified,
		"userType":           aliRsp.Response.UserType,
		"userState":          aliRsp.Response.UserStatus,
		"isCertified":        aliRsp.Response.IsCertified,
		"sex":                aliRsp.Response.Gender,
		"authToken":          token.Response.AccessToken,
	}
	// 存起来，存到缓存即可，redis  直接存到数据库
	//keyId := aliRsp.Response.UserId
	//// 设置缓存
	//gcache.Set(ctx, keyId, authInfo, time.Hour)
	consumerInfo := alipay_model.AlipayConsumerConfig{}
	gconv.Struct(authInfo, &consumerInfo)
	_, err = service.Consumer().CreateConsumer(ctx, consumerInfo)
	if err != nil {
		return false
	}

	return true
}
