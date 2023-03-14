package app_auth

import (
    "context"
    "fmt"
    "github.com/go-pay/gopay"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
    "github.com/gogf/gf/v2/util/gconv"
    "github.com/kysion/alipay-test/alipay_model"
    enum "github.com/kysion/alipay-test/alipay_model/alipay_enum"
    service "github.com/kysion/alipay-test/alipay_service"
    "github.com/kysion/alipay-test/internal/logic/internal/aliyun"
)

/*
	授权后将商户AppAuthToken存储起来
*/

type sAppAuth struct {
}

func NewAppAuth() *sAppAuth {
    // 初始化文件内容

    result := &sAppAuth{}

    result.injectHook()
    return result
}

func (s *sAppAuth) injectHook() {
    service.Gateway().InstallHook(enum.Info.Type.AlipayAppAuth, s.AppAuth)
}

// AppAuth 具体服务
func (s *sAppAuth) AppAuth(ctx context.Context, info g.Map) bool {
    fmt.Println("hello authApp")

    client, _ := aliyun.NewClient(ctx, "")
    data := gopay.BodyMap{}
    gconv.Struct(info, &data)
    // data.Set("code", data.Get("app_auth_code")) // 商家授权code
    // 根据商家code获取到token 需要存储下来
    aliRsp, _ := client.OpenAuthTokenApp(ctx, data)
    fmt.Println(aliRsp)

    // 存起来
    authInfos := g.Map{}

    if len(aliRsp.Response.Tokens) != 0 {
        token := aliRsp.Response.Tokens[0]
        authInfos = g.Map{
            "app_auth_token":    token.AppAuthToken,
            "app_refresh_token": token.AppRefreshToken,
            "auth_app_id":       token.AuthAppId,
            "expires_in":        token.ExpiresIn,
            "re_expires_in":     token.ReExpiresIn,
            "user_id":           token.UserId,
        }
    } else {
        authInfos = g.Map{
            "app_auth_token":    aliRsp.Response.AppAuthToken,
            "app_refresh_token": aliRsp.Response.AppRefreshToken,
            "auth_app_id":       aliRsp.Response.AuthAppId,
            "expires_in":        aliRsp.Response.ExpiresIn,
            "re_expires_in":     aliRsp.Response.ReExpiresIn,
            "user_id":           aliRsp.Response.UserId,
        }
    }

    // 存起来，存到缓存即可，redis
    // keyId := gconv.String(idgen.NextId())
    //keyId := aliRsp.Response.UserId
    //
    //gcache.Set(ctx, keyId, authInfos, time.Hour)

    // 商家Token更新 (一般是商家在我们系统中有了商家之后，我们进行修改企App_auth_token)
    config := alipay_model.UpdateMerchantAppAuthToken{
        AppId:        gconv.String(authInfos["auth_app_id"]),
        AppAuthToken: gconv.String(authInfos["app_auth_token"]),
        ExpiresIn:    gtime.New(authInfos["expires_in"]),
        ReExpiresIn:  gtime.New(authInfos["re_expires_in"]),
        ThirdAppId:   data.Get("app_id"),
    }

    app, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, config.AppId)

    if err == nil && app != nil {
        _, err := service.MerchantAppConfig().UpdateAppAuthToken(ctx, &config)
        if err != nil {
            return false
        }
    } else if app == nil {
        // 如果暂时
        //config := alipay_model.AlipayMerchantAppConfig{}
        //gconv.Struct(authInfos, &config)
        //_, err := service.MerchantAppConfig().CreateMerchantAppConfig(ctx, &config)
        //if err != nil {
        //    return false
        //}
    }

    return true
}
