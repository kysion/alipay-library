package merchant

import (
	"context"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/go-pay/gopay"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
)

/*
	授权后将商户AppAuthToken存储起来
*/

type sAppAuth struct {
}

func NewAppAuth() service.IAppAuth {
	// 初始化文件内容

	result := &sAppAuth{}

	result.injectHook()
	return result
}

func (s *sAppAuth) injectHook() {
	hook := service.Gateway().GetCallbackMsgHook()

	hook.InstallHook(enum.Info.CallbackType.AlipayAppAuth, s.AppAuth)
}

// AppAuth 具体服务
func (s *sAppAuth) AppAuth(ctx context.Context, info g.Map) string { // 返回商户userId
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay应用授权 ------- ", "sAppAuth")

	fmt.Println("hello authApp")

	data := gopay.BodyMap{}
	gconv.Struct(info, &data)

	if data.Get("app_id") == "" {
		return ""
	}

	client, _ := aliyun.NewClient(ctx, data.Get("app_id"))

	// data.Set("code", data.Get("app_auth_code")) // 商家授权code

	// 1.根据商家code获取到token 需要存储下来
	aliRsp, _ := client.OpenAuthTokenApp(ctx, data)
	fmt.Println(aliRsp)

	// 2.存起来
	gconv.Int64(data.Get("sys_user_id")) // 授权码的附带数据sys_user_id

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

	err := dao.AlipayMerchantAppConfig.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {

		// 商家Token更新 (一般是商家在我们系统中有了商家之后，我们进行修改企App_auth_token)  后续业务才有，筷满客只需要记录商家唯一标识
		config := alipay_model.UpdateMerchantAppAuthToken{
			AppId:        gconv.String(authInfos["auth_app_id"]),
			AppAuthToken: gconv.String(authInfos["app_auth_token"]),
			ExpiresIn:    gtime.New(authInfos["expires_in"]),
			ReExpiresIn:  gtime.New(authInfos["re_expires_in"]),
			ThirdAppId:   data.Get("app_id"),
		}

		app, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, config.AppId)

		if err == nil && app != nil { // 有商家应用的，就可以直接修改AppAuthToken，并且添加一条第三方和用户关系记录
			_, err := service.MerchantAppConfig().UpdateAppAuthToken(ctx, &config)
			if err != nil {
				return err
			}
		} else if app == nil { // 如果没有的，添加一条商家应用配置和用户关系记录???????????????
			//merchantConfig := alipay_model.AlipayMerchantAppConfig{
			//
			//}
			//
			//gconv.Struct(authInfos, &config)
			//_, err := service.MerchantAppConfig().CreateMerchantAppConfig(ctx, &merchantConfig)
			//if err != nil {
			//   return false
			//}
		}

		// 3.记录第三方应用和商家用户关系  (Hook创建)
		//platformUser, err := share_service.PlatformUser().GetPlatFormUserByUserAndMerchantAppId(ctx, gconv.String(authInfos["user_id"]), config.AppId)
		//
		//platformUserId := gconv.String(authInfos["user_id"])
		//
		//if err == nil && platformUser != nil { // 说明存在
		//    // 进行Update
		//    platform := share_model.UpdatePlatformUser{
		//        Id:         platformUser.Id,
		//        EmployeeId: sysUserId,
		//    }
		//    if authInfos["user_id"] != "" {
		//        platform.UserId = platformUserId
		//    }
		//
		//    rows, err := share_service.PlatformUser().UpdatePlatformUser(ctx, &platform)
		//    if err != nil || rows == nil {
		//        return err
		//    }
		//
		//} else if platformUser == nil { // 不存在创建
		//    // 根据sys_user_id查询商户信息
		//    employee, err := share_consts.Global.Merchant.Employee().GetEmployeeById(ctx, sysUserId)
		//
		//    platform := share_model.PlatformUser{
		//        Id:            idgen.NextId(),
		//        FacilitatorId: 0,
		//        OperatorId:    0,
		//        EmployeeId:    sysUserId,
		//        MerchantId:    0,
		//        Platform:      pay_enum.Order.TradeSourceType.Alipay.Code(), // 来源
		//        ThirdAppId:    gconv.String(data.Get("app_id")),
		//        MerchantAppId: gconv.String(authInfos["auth_app_id"]),
		//        UserId:        platformUserId,
		//        Type:          share_enum.User.Type.Merchant.Code(), // 用户类型商户
		//    }
		//
		//    if employee != nil {
		//        platform.EmployeeId = employee.Id
		//        platform.MerchantId = employee.UnionMainId
		//    }
		//
		//    rows, err := share_service.PlatformUser().CreatePlatformUser(ctx, &platform)
		//    if err != nil || rows == nil {
		//        return err
		//    }
		//}
		return nil
	})

	if err != nil {
		return ""
	}

	// 4.写接口，根据平台用户id获取JwtToken
	//g.RequestFromCtx(ctx).Response.RedirectTo("")

	// 返回商户userId
	return gconv.String(authInfos["user_id"])
}
