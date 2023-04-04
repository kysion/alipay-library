package merchant

import (
	"context"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/sys_model/sys_enum"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/gopay"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/yitter/idgenerator-go/idgen"
)

/*
	获取支付宝会员信息等
*/

type sWallet struct {
	ConsumerHook base_hook.BaseHook[enum.ConsumerAction, hook.ConsumerHookFunc]
}

func NewWallet() *sWallet {
	// 初始化文件内容

	result := &sWallet{}

	result.injectHook()
	fmt.Println(result)
	return result
}

func (s *sWallet) injectHook() {
	hook := service.Gateway().GetCallbackMsgHook()

	hook.InstallHook(enum.Info.CallbackType.AlipayWallet, s.Wallet)
}

func (s *sWallet) InstallConsumerHook(infoType enum.ConsumerAction, hookFunc hook.ConsumerHookFunc) {
	s.ConsumerHook.InstallHook(infoType, hookFunc)
}

// Wallet 具体服务 H5用户授权 + 小程序
func (s *sWallet) Wallet(ctx context.Context, info g.Map) string {
	res := alipay_model.UserInfoShare{}

	client, err := aliyun.NewClient(ctx, gconv.String(info["app_id"]))

	data := gopay.BodyMap{}
	gconv.Struct(info, &data)

	err = dao.AlipayMerchantAppConfig.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 根据AppId获取商家相关配置，包括AppAuthToken
		merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(info["app_id"]))
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
		sys_service.SysUser().MakeSession(ctx, merchantApp.SysUserId)

		//employee, err := share_consts.Global.Merchant.Employee().GetEmployeeById(ctx, merchantApp.SysUserId)

		// 3.存储消费者数据并创建用户  kmk-consumer
		var consumerId int64 // 消费者Id

		if token.Response.UserId != "" { // 发布授权广播
			s.ConsumerHook.Iterator(func(key enum.ConsumerAction, value hook.ConsumerHookFunc) {
				if key.Code() == enum.Consumer.ActionEnum.Auth.Code() { // 如果订阅者是订阅授权
					g.Try(ctx, func(ctx context.Context) {
						data := hook.UserInfo{
							SysUserId:     merchantApp.SysUserId,
							UserInfoShare: *userInfo.Response,
						}
						// 返回消费者id
						consumerId = value(ctx, data)
					})
				}
			})
		}

		//getConsumer, err := share_service.ConsumerConfig().GetConsumerByAlipayUnionId(ctx, userInfo.Response.UserId)
		//consumerRes := &share_model.ConsumerRes{}
		//if err != nil && getConsumer == nil {
		//    shareConsumer := kconv.Struct(userInfo.Response, &share_model.Consumer{})
		//    shareConsumer.AlipayUnionId = gconv.String(userInfo.Response.UserId)
		//
		//    if employee != nil {
		//        shareConsumer.SysUserId = gconv.String(employee.Id)
		//    }
		//
		//    consumerRes, err = share_service.ConsumerConfig().CreateConsumer(ctx, shareConsumer)
		//    if err != nil {
		//        return err
		//    }
		//} else {
		//    //consumerInfo := share_model.UpdateConsumer{}
		//    //gconv.Struct(userInfo.Response, &consumerInfo)
		//    //consumerInfo.Id = merchantApp.SysUserId
		//    //
		//    //_, err = share_service.ConsumerConfig().UpdateConsumer(ctx, &consumerInfo)
		//    //if err != nil {
		//    //	return err
		//    //}
		//}

		// 4.存储阿里消费者记录 alipay-consumer-config
		alipayConsumer, err := service.ConsumerConfig().GetConsumerByUserId(ctx, userInfo.Response.UserId)

		if err != nil && alipayConsumer == nil { // 消费者不存在，则创建
			consumerInfo := alipay_model.AlipayConsumerConfig{}
			gconv.Struct(userInfo.Response, &consumerInfo)
			consumerInfo.UserId = gconv.String(userInfo.Response.UserId)

			//if employee != nil {
			//    consumerInfo.SysUserId = employee.Id
			//}

			if consumerId != 0 {
				consumerInfo.SysUserId = consumerId
			}

			_, err = service.ConsumerConfig().CreateConsumer(ctx, &consumerInfo)
			if err != nil {
				return err
			}

		} else { // 存在则更新
			consumerInfo := alipay_model.UpdateConsumerReq{}
			gconv.Struct(userInfo.Response, &consumerInfo)

			_, err = service.ConsumerConfig().UpdateConsumer(ctx, alipayConsumer.Id, &consumerInfo)
			if err != nil {
				return err
			}
		}

		// 5.存储第三方应用和用户关系记录
		s.ConsumerHook.Iterator(func(key enum.ConsumerAction, value hook.ConsumerHookFunc) {
			if key.Code() == enum.Consumer.ActionEnum.Auth.Code() { // 如果订阅者是订阅授权
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
						UserId:        token.Response.UserId,                // 平台账户唯一标识
						Type:          sys_enum.User.Type.New(0, "").Code(), // 用户类型匿名消费者
					}
					//
					//if employee != nil {
					//    data.EmployeeId = employee.Id
					//    data.MerchantId = employee.UnionMainId
					//}

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
		return ""
	}

	// 返回用户在平台的唯一id
	return res.UserId
}
