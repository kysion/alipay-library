package merchant_controller

import (
	"context"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/api_v1"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/api/alipay_v1/alipay_merchant_v1"
	"strconv"
)

var MerchantService = cMerchantService{}

type cMerchantService struct{}

// AuthMerchantApp 商家授权
func (c *cMerchantService) AuthMerchantApp(ctx context.Context, _ *alipay_merchant_v1.AuthMerchantAppReq) (api_v1.StringRes, error) {
	info := g.RequestFromCtx(ctx).Get("appId").String()

	appId, _ := strconv.ParseInt(info, 32, 0)

	app, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(appId))
	if err != nil {
		return "", err
	}

	data := g.Map{
		"platformCode": "O",
		"taskType":     "INTERFACE_AUTH",
		"agentOpParam": g.Map{
			//"redirectUri": "https//alipay.kuaimk.com/alipay/1975251903f95826/gateway.callback?sys_user_id=6343708512288837",
			"redirectUri": app.AppCallbackUrl,
			"appTypes": []string{
				"TINYAPP",
				"WEBAPP",
				"PUBLICAPP",
			},
			"isvAppId": app.ThirdAppId,
		},
	}

	encode, _ := gjson.Encode(data)
	params := gurl.Encode("pages/authorize/index?bizData=" + string(encode))
	url := "alipays://platformapi/startapp?appId=2021003130652097&page=" + params

	g.RequestFromCtx(ctx).Response.RedirectTo(url)
	//g.RequestFromCtx(ctx).Response.RedirectTo("pages%2Fauthorize%2Findex%3FbizData%3D%7B%22platformCode%22%3A%22O%22%2C%22taskType%22%3A%22INTERFACE_AUTH%22%2C%22agentOpParam%22%3A%7B%22redirectUri%22%3A%22https%3A%2F%2Falipay.kuaimk.com%2Falipay%2F1975251903f95826%2Fgateway.callback%3Fsys_user_id%3D6343708512288837%22%2C%22appTypes%22%3A%5B%22TINYAPP%22%2C%22WEBAPP%22%2C%22PUBLICAPP%22%2C%22MOBILEAPP%22%5D%2C%22isvAppId%22%3A%222021003179681073%22%7D%7D")

	return "", nil
}

// GetAlipayUserInfo 获取支付宝会员信息，相当于静默登录  H5
func (c *cMerchantService) GetAlipayUserInfo(ctx context.Context, _ *alipay_merchant_v1.GetAlipayUserInfoReq) (api_v1.StringRes, error) {
	info := g.RequestFromCtx(ctx).Get("appId").String()

	appId, _ := strconv.ParseInt(info, 32, 0)

	app, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(appId))
	if err != nil {
		return "", err
	}
	url := "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm?app_id=" + app.AppId + "&scope=auth_user&redirect_uri="

	//	url += gurl.Encode(app.AppCallbackUrl + "?sys_user_id=" + gconv.String(app.SysUserId) + "&merchant_id=" + gconv.String(app.UnionMainId))

	url += gurl.Encode(app.AppCallbackUrl)

	//url += gurl.Encode("https://alipay.kuaimk.com/alipay/1pe2ts9qp9h/gateway.callback")

	fmt.Println(url)

	g.RequestFromCtx(ctx).Response.RedirectTo(url)

	// 网页移动应用可以获取用户信息：   回调地址需要换成kuaimk的  https://alipay.kuaimk.com/alipay/1975251903f95826/gateway.callback
	// g.RequestFromCtx(ctx).Response.RedirectTo("https%3A%2F%2Falipay.kuaimk.com%2Falipay%2F1975251903f95826%2Fgateway.callback")

	// 三方应用和小程序不能获取用户信息，原因如下：
	//     - 小程序应用没有授权回调地址，再换个网页或者、移动应用
	//     - 三方应用是在三方调用情况下使用的，您现在是自调用模式调用 merchant.user.info.share（支付宝会员授权信息查询接口）会因为不是三方调用报此用户不允许自调用

	return "", nil
}

// GetUserInfoByAuthCode 根据认证码获取会员信息并注册  认证码前端提供 小程序
//func (c *cMerchantService) GetUserInfoByAuthCode(ctx context.Context, req *alipay_merchant_v1.GetUserInfoByAuthCodeReq) (*alipay_model.UserInfoRes, error) {
//    ret, err := service.MerchantService().UserInfoAuth(ctx, req.AuthCode, req.AppId,)
//
//    return (*alipay_model.UserInfoRes)(ret), err
//}
