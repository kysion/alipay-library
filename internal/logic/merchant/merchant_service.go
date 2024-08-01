package merchant

import (
	"context"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/SupenBysz/gf-admin-community/utility/idgen"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
	entity "github.com/kysion/alipay-library/alipay_model/alipay_entity"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/gopay"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
)

/*
	商户相关API服务
*/

type sMerchantService struct {
	// 消费者Hook
	ConsumerHook base_hook.BaseHook[hook.ConsumerKey, hook.ConsumerHookFunc]

	// 平台与用户关联Hook
	//PlatFormUserHook base_hook.BaseHook[alipay_enum.ConsumerAction, hook.PlatFormUserHookFunc]
}

func NewMerchantService() *sMerchantService {
	result := &sMerchantService{}

	result.injectHook()

	return result
}

func (s *sMerchantService) injectHook() {
	hook := service.Gateway().GetCallbackMsgHook()

	hook.InstallHook(alipay_enum.Info.CallbackType.AlipayWallet, s.UserInfoAuth)
}

func (s *sMerchantService) InstallConsumerHook(infoType hook.ConsumerKey, hookFunc hook.ConsumerHookFunc) {
	sys_service.SysLogs().InfoSimple(context.Background(), nil, "\n-------订阅sMerchantService-Hook： ------- ", "sPlatformUser")

	s.ConsumerHook.InstallHook(infoType, hookFunc)
}

func (s *sMerchantService) GetHook() base_hook.BaseHook[hook.ConsumerKey, hook.ConsumerHookFunc] {
	return s.ConsumerHook
}

// GetUserId 用于检查是否注册,如果已经注册，返会userId
func (s *sMerchantService) GetUserId(ctx context.Context, authCode string, appId string) (response *alipay.SystemOauthTokenResponse, err error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------获取用户JwtToken----GetUserId---- ", "sMerchantService")

	//userId := ""
	//err = dao.AlipayMerchantAppConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
	var client *aliyun.AliPay

	// 根据AppId获取商家相关配置，包括AppAuthToken
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil || merchantApp == nil {
		return nil, err
	}

	if merchantApp.ThirdAppId != "" {
		client, err = aliyun.NewClient(ctx, appId)
	} else {
		client, err = aliyun.NewMerchantClient(ctx, appId)
	}

	client.SetAppAuthToken(merchantApp.AppAuthToken)

	// 1.auth_code换Token
	token, _ := client.SystemOauthToken(ctx, gopay.BodyMap{
		"code":       authCode,
		"grant_type": "authorization_code",
	})

	//fmt.Println("平台用户id：", token.Response.UserId)
	//userId = token.Response.UserId

	//return "", nil
	//})
	//
	//if err != nil {
	//	return "", err
	//}

	return token, nil
}

// UserInfoAuth 具体服务 用户授权 + 小程序和H5都兼容
func (s *sMerchantService) UserInfoAuth(ctx context.Context, info g.Map) string { // code string, appId string, sysUserId ...int64
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------支付宝用户授权 UserInfoAuth", "sMerchantService")

	from := gmap.NewStrAnyMapFrom(info)
	var client *aliyun.AliPay
	code := from.Get("code") //
	fmt.Println(code)
	appId := gconv.String(from.Get("app_id"))
	sysUserId := gconv.Int64(from.Get("sys_user_id"))
	//sysUserId = 8138722065186885
	merchantId := gconv.Int64(from.Get("merchant_id"))
	tokenStr := from.Get("token")
	token := kconv.Struct(tokenStr, &alipay.SystemOauthTokenResponse{})

	res := alipay_model.UserInfoShare{}
	// 根据AppId获取商家相关配置，包括AppAuthToken
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil || merchantApp == nil {
		return ""
	}

	if merchantApp.ThirdAppId != "" {
		client, err = aliyun.NewClient(ctx, appId)
	} else {
		client, err = aliyun.NewMerchantClient(ctx, appId)
	}

	data := gopay.BodyMap{}
	gconv.Struct(info, &data)

	err = dao.AlipayMerchantAppConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {

		// 这个token是动态的，哪个商家需要获取，appId和appAuthToken就传递对应的
		if merchantApp.AppAuthToken != "" {
			client.SetAppAuthToken(merchantApp.AppAuthToken) // 商家Token
		}

		if token.Response == nil { // 避免用过一次的情况
			// 1.auth_code换access_token
			token, _ = client.SystemOauthToken(ctx, data)
			fmt.Println(token)
			fmt.Println("平台用户id：", token.Response.UserId)
			res.UserId = token.Response.UserId
		}

		// 2.token获取支付宝会员授权信息查询接口  小程序好像查不到
		aliRsp, _ := client.UserInfoShare(ctx, token.Response.AccessToken)
		fmt.Println(aliRsp)

		userInfo := alipay_model.UserInfoShareResponse{}
		gconv.Struct(aliRsp, &userInfo)
		// 将返回值赋值
		gconv.Struct(userInfo.Response, &res) // 小程序的静默授权拿不到userInfo，只能拿到userId
		userInfo.Response.UserId = token.Response.UserId

		// 根据sys_user_id查询商户信息
		sys_service.SysUser().MakeSession(ctx, sysUserId)
		sysUser, err := sys_service.SysUser().GetSysUserById(ctx, sysUserId)
		if err != nil && sysUser == nil {
			// TODO 后续需要删除下面创建用户的代码。需要在业务层先创建好用户
			// 1.创建一个匿名sysUser = 用户登录信息
			/*id := idgen.NextId()
			passwordLen := len(gconv.String(id))
			pwd := gstr.SubStr(gconv.String(id), passwordLen-6, 6) // 密码为id后六位

			data := sys_model.UserInnerRegister{
				Username:        strconv.FormatInt(gconv.Int64(id), 36),
				Password:        pwd,
				ConfirmPassword: pwd,
			}
			sysUser, err = sys_service.SysUser().CreateUser(ctx, data, sys_enum.User.State.Normal, sys_enum.User.Type.New(0, "匿名"), id)
			if err != nil {
				return err
			}
			sysUserId = sysUser.Id*/

			return err
		}

		//employee, err := share_consts.Global.Merchant.Employee().GetEmployeeById(ctx,sysUserId)

		// 3.存储消费者数据   kmk-consumer
		if token.Response.UserId != "" { // 发布授权广播
			s.ConsumerHook.Iterator(func(key hook.ConsumerKey, value hook.ConsumerHookFunc) {
				if key.ConsumerAction.Code() == alipay_enum.Consumer.ActionEnum.Auth.Code() && key.Category.Code() == alipay_enum.Consumer.Category.Consumer.Code() { // 如果订阅者是订阅授权,并且是操作kmk_consumer表
					g.Try(ctx, func(ctx context.Context) {
						data := hook.UserInfo{
							SysUserId:     sysUser.Id, // (消费者id = sys_User_id)
							UserInfoShare: *userInfo.Response,
						}

						sys_service.SysLogs().InfoSimple(ctx, nil, "\n广播-------存储消费者数据 kmk-consumer", "sMerchantService")

						//consumerId = value(ctx, data)
						value(ctx, data)
					})
				}
			})
		}

		// 4.存储阿里消费者记录 alipay-consumer-config
		//alipayConsumer, _ := service.ConsumerConfig().GetConsumerByUserId(ctx, userInfo.Response.UserId)
		alipayConsumer, _ := service.ConsumerConfig().GetConsumerByUserIdAndAppId(ctx, userInfo.Response.UserId, appId)

		if alipayConsumer == nil { // 消费者不存在，则创建
			consumerInfo := alipay_model.AlipayConsumerConfig{}
			gconv.Struct(userInfo.Response, &consumerInfo)
			consumerInfo.AppId = appId
			consumerInfo.AppType = merchantApp.AppType

			consumerInfo.UserId = gconv.String(userInfo.Response.UserId)
			consumerInfo.AuthToken = token.Response.AccessToken
			consumerInfo.ReFreshToken = token.Response.RefreshToken
			consumerInfo.AuthState = alipay_enum.Consumer.AuthState.Auth.Code()
			consumerInfo.AuthStart = gtime.NewFromStr(token.Response.AuthStart)
			consumerInfo.AlipayUserId = token.Response.AlipayUserId
			consumerInfo.ExpiresIn = gtime.NewFromTimeStamp(gtime.Now().Timestamp() + gconv.Int64(token.Response.ExpiresIn))
			consumerInfo.ReExpiresIn = gtime.NewFromTimeStamp(gtime.Now().Timestamp() + gconv.Int64(token.Response.ReExpiresIn))

			//consumerInfo. =  token.Response.AccessToken
			consumerInfo.SysUserId = sysUser.Id

			_, err = service.ConsumerConfig().CreateConsumer(ctx, &consumerInfo)
			if err != nil {
				return err
			}

		} else { // 存在则更新
			consumerInfo := alipay_model.UpdateConsumerReq{}
			authEnum := alipay_enum.Consumer.AuthState.Auth.Code()
			gconv.Struct(userInfo.Response, &consumerInfo)
			consumerInfo.AuthToken = &token.Response.AccessToken
			consumerInfo.ReFreshToken = &token.Response.RefreshToken
			consumerInfo.AuthState = &authEnum
			consumerInfo.AuthStart = gtime.NewFromStr(token.Response.AuthStart)
			consumerInfo.AlipayUserId = &token.Response.AlipayUserId
			consumerInfo.ExpiresIn = gtime.NewFromTimeStamp(gtime.Now().Timestamp() + gconv.Int64(token.Response.ExpiresIn))
			consumerInfo.ReExpiresIn = gtime.NewFromTimeStamp(gtime.Now().Timestamp() + gconv.Int64(token.Response.ReExpiresIn))

			_, err = service.ConsumerConfig().UpdateConsumer(ctx, alipayConsumer.Id, &consumerInfo)
			if err != nil {
				return err
			}
		}

		// 5.存储第三方应用和用户关系记录
		s.ConsumerHook.Iterator(func(key hook.ConsumerKey, value hook.ConsumerHookFunc) { // 这会同时走两个Hook，kmk_consumer  + platform_user
			if key.ConsumerAction.Code() == alipay_enum.Consumer.ActionEnum.Auth.Code() && key.Category.Code() == alipay_enum.Consumer.Category.PlatFormUser.Code() { // 如果订阅者是订阅授权
				g.Try(ctx, func(ctx context.Context) {
					platformUser := entity.PlatformUser{
						Id:             idgen.NextId(),
						FacilitatorId:  0,
						OperatorId:     0,
						SysUserId:      sysUser.Id,                                   // EmployeeId  == consumerId == sysUserId   三者相等
						MerchantId:     merchantId,                                   // 商家id，就是消费者首次扫码的商家
						PlatformType:   pay_enum.Order.TradeSourceType.Alipay.Code(), // 平台类型：1支付宝、2微信、4抖音、8银联
						ThirdAppId:     merchantApp.ThirdAppId,
						MerchantAppId:  merchantApp.AppId,
						PlatformUserId: token.Response.UserId, // 平台账户唯一标识
						SysUserType:    sysUser.Type,          // 用户类型
					}

					sys_service.SysLogs().InfoSimple(ctx, nil, "\n广播-------存储第三方应用和用户关系记录 kmk-plat_form_user", "sMerchantService")

					value(ctx, platformUser) // 调用Hook
				})
			}
		})

		return nil
	})

	if err != nil {
		return ""
	}

	// 返回用户在阿里的唯一标识userId
	return res.UserId
}
