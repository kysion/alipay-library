package merchant

import (
	"context"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/go-pay/gopay"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kuaimk/kmk-share-library/share_model/share_enum"
	"github.com/kysion/alipay-library/alipay_model"
	dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/yitter/idgenerator-go/idgen"
)

/*
	商户相关API服务
*/

type sMerchantService struct {
	// 消费者Hook
	ConsumerHook base_hook.BaseHook[alipay_enum.ConsumerAction, hook.ConsumerHookFunc]

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

func (s *sMerchantService) InstallConsumerHook(infoType alipay_enum.ConsumerAction, hookFunc hook.ConsumerHookFunc) {
	s.ConsumerHook.InstallHook(infoType, hookFunc)
	fmt.Println(s.ConsumerHook)
}

func (s *sMerchantService) GetHook() base_hook.BaseHook[alipay_enum.ConsumerAction, hook.ConsumerHookFunc] {
	return s.ConsumerHook
}

// GetUserId 用于检查是否注册,如果已经注册，返会userId
func (s *sMerchantService) GetUserId(ctx context.Context, authCode string, appId string) (res string, err error) {
	client, err := aliyun.NewClient(ctx, appId)
	userId := ""
	err = dao.AlipayMerchantAppConfig.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {

		// 根据AppId获取商家相关配置，包括AppAuthToken
		merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
		if err != nil || merchantApp == nil {
			return err
		}

		client.SetAppAuthToken(merchantApp.AppAuthToken)

		// 1.auth_code换Token
		token, _ := client.SystemOauthToken(ctx, gopay.BodyMap{
			"code":       authCode,
			"grant_type": "authorization_code",
		})

		fmt.Println("平台用户id：", token.Response.UserId)
		userId = token.Response.UserId

		return nil
	})

	if err != nil {
		return "", err
	}

	return userId, nil
}

// UserInfoAuth 具体服务 用户授权 + 小程序和H5都兼容
func (s *sMerchantService) UserInfoAuth(ctx context.Context, info g.Map) bool { // code string, appId string, sysUserId ...int64
	from := gmap.NewStrAnyMapFrom(info)

	//code := from.Get("code") //
	appId := gconv.String(from.Get("app_id"))
	sysUserId := gconv.Int64(from.Get("sys_user_id"))
	merchantId := gconv.Int64(from.Get("merchant_id"))

	res := alipay_model.UserInfoShare{}

	client, err := aliyun.NewClient(ctx, appId)

	data := gopay.BodyMap{}
	gconv.Struct(info, &data)

	err = dao.AlipayMerchantAppConfig.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 根据AppId获取商家相关配置，包括AppAuthToken
		merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
		if err != nil || merchantApp == nil {
			return err
		}

		// data.Set("code", data.Get("auth_code"))                            // 用户授权code

		// 这个token是动态的，哪个商家需要获取，appId和appAuthToken就传递对应的
		client.SetAppAuthToken(merchantApp.AppAuthToken) // 商家Token

		// 1.auth_code换Token
		token, _ := client.SystemOauthToken(ctx, data)

		fmt.Println("平台用户id：", token.Response.UserId)

		// 2.token获取支付宝会员授权信息查询接口
		aliRsp, _ := client.UserInfoShare(ctx, token.Response.AccessToken)
		fmt.Println(token)
		fmt.Println(aliRsp)

		userInfo := alipay_model.UserInfoShareResponse{}
		gconv.Struct(aliRsp, &userInfo)
		// 将返回值赋值
		gconv.Struct(userInfo.Response, &res) // 小程序静默授权拿不到userInfo，只能拿到userId
		res.UserId = token.Response.UserId
		userInfo.Response.UserId = token.Response.UserId

		// 根据sys_user_id查询商户信息
		sys_service.SysUser().MakeSession(ctx, sysUserId)

		user, err := sys_service.SysUser().GetSysUserById(ctx, sysUserId)
		if err != nil {
			return err
		}

		//employee, err := share_consts.Global.Merchant.Employee().GetEmployeeById(ctx,sysUserId)

		// 3.存储消费者数据并创建用户  kmk-consumer
		var consumerId int64 // 消费者Id

		if token.Response.UserId != "" { // 发布授权广播
			s.ConsumerHook.Iterator(func(key alipay_enum.ConsumerAction, value hook.ConsumerHookFunc) {
				if key.Code() == alipay_enum.Consumer.ActionEnum.Auth.Code() { // 如果订阅者是订阅授权
					g.Try(ctx, func(ctx context.Context) {
						data := hook.UserInfo{
							SysUserId:     user.Id,
							UserInfoShare: *userInfo.Response,
						}
						consumerId = value(ctx, data)
					})
				}
			})
		}

		//getConsumer, err := share_service.Consumer().GetConsumerByAlipayUnionId(ctx, userInfo.Response.UserId)
		//consumerRes := &share_model.ConsumerRes{}
		//if err != nil && getConsumer == nil {
		//    shareConsumer := kconv.Struct(userInfo.Response, &share_model.Consumer{})
		//    shareConsumer.AlipayUnionId = gconv.String(userInfo.Response.UserId)
		//
		//    if employee != nil {
		//        shareConsumer.SysUserId = gconv.String(employee.Id)
		//    }
		//
		//    consumerRes, err = share_service.Consumer().CreateConsumer(ctx, shareConsumer)
		//    if err != nil {
		//        return err
		//    }
		//} else {
		//    //consumerInfo := share_model.UpdateConsumer{}
		//    //gconv.Struct(userInfo.Response, &consumerInfo)
		//    //consumerInfo.Id = merchantApp.SysUserId
		//    //
		//    //_, err = share_service.Consumer().UpdateConsumer(ctx, &consumerInfo)
		//    //if err != nil {
		//    //	return err
		//    //}
		//}

		// 4.存储阿里消费者记录 alipay-consumer-config
		alipayConsumer, err := service.Consumer().GetConsumerByUserId(ctx, userInfo.Response.UserId)

		if err != nil && alipayConsumer == nil { // 消费者不存在，则创建
			consumerInfo := alipay_model.AlipayConsumerConfig{}
			gconv.Struct(userInfo.Response, &consumerInfo)
			consumerInfo.UserId = gconv.String(userInfo.Response.UserId)

			if consumerId != 0 {
				consumerInfo.SysUserId = consumerId
			}

			_, err = service.Consumer().CreateConsumer(ctx, consumerInfo)
			if err != nil {
				return err
			}

		} else { // 存在则更新
			consumerInfo := alipay_model.UpdateConsumerReq{}
			gconv.Struct(userInfo.Response, &consumerInfo)

			_, err = service.Consumer().UpdateConsumer(ctx, alipayConsumer.Id, consumerInfo)
			if err != nil {
				return err
			}
		}

		// 5.存储第三方应用和用户关系记录
		s.ConsumerHook.Iterator(func(key alipay_enum.ConsumerAction, value hook.ConsumerHookFunc) {
			if key.Code() == alipay_enum.Consumer.ActionEnum.Auth.Code() { // 如果订阅者是订阅授权
				g.Try(ctx, func(ctx context.Context) {
					data := hook.PlatformUser{
						Id:            idgen.NextId(),
						FacilitatorId: 0,
						OperatorId:    0,
						EmployeeId:    consumerId,
						MerchantId:    0,
						Platform:      pay_enum.Order.TradeSourceType.Alipay.Code(), // 来源
						ThirdAppId:    merchantApp.ThirdAppId,
						MerchantAppId: merchantApp.AppId,
						UserId:        token.Response.UserId,                 // 平台账户唯一标识
						Type:          share_enum.User.Type.Anonymous.Code(), // 用户类型匿名消费者
					}

					if user != nil {
						data.EmployeeId = user.Id
					}

					data.MerchantId = merchantId

					if consumerId != 0 { // 适用于消费者没有员工的情况下
						data.EmployeeId = consumerId
					}

					value(ctx, data) // 调用Hook
				})
			}
		})
		//platUser, err := share_service.PlatformUser().GetPlatformUserByUserId(ctx, userInfo.Response.UserId)
		//
		//if err != nil && platUser == nil { // 不存在创建
		//    platform := share_model.PlatformUser{
		//        Id:            idgen.NextId(),
		//        FacilitatorId: 0,
		//        OperatorId:    0,
		//        EmployeeId:    consumerRes.Id,
		//        MerchantId:    0,
		//        Platform:      pay_enum.Order.TradeSourceType.Alipay.Code(), // 来源
		//        ThirdAppId:    merchantApp.ThirdAppId,
		//        MerchantAppId: merchantApp.AppId,
		//        UserId:        userInfo.Response.UserId,              // 平台账户唯一标识
		//        Type:          share_enum.User.Type.Anonymous.Code(), // 用户类型匿名消费者
		//    }
		//
		//    if employee != nil {
		//        platform.EmployeeId = employee.Id
		//        platform.MerchantId = employee.UnionMainId
		//    }
		//
		//    if consumerRes != nil && consumerRes.Id != 0 { // 适用于消费者没有员工的情况下
		//        platform.EmployeeId = consumerRes.Id
		//    }
		//
		//    _, err = share_service.PlatformUser().CreatePlatformUser(ctx, &platform)
		//    if err != nil {
		//        return err
		//    }
		//
		//}

		return nil
	})

	if err != nil {
		return false
	}

	return true
}
