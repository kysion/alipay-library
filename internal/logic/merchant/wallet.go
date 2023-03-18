package merchant

import (
    "context"
    "fmt"
    "github.com/go-pay/gopay"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/util/gconv"
    "github.com/kuaimk/kmk-share-library/share_consts"
    "github.com/kuaimk/kmk-share-library/share_model"
    "github.com/kuaimk/kmk-share-library/share_model/share_enum"
    "github.com/kuaimk/kmk-share-library/share_service"
    "github.com/kysion/alipay-library/alipay_model"
    dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
    enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
    service "github.com/kysion/alipay-library/alipay_service"
    "github.com/kysion/alipay-library/internal/logic/internal/aliyun"
    "github.com/kysion/base-library/utility/kconv"
    "github.com/kysion/pay-share-library/pay_model/pay_enum"
    "github.com/yitter/idgenerator-go/idgen"
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

// Wallet 具体服务 H5用户授权
func (s *sWallet) Wallet(ctx context.Context, info g.Map) bool {
    client, _ := aliyun.NewClient(ctx, "")

    data := gopay.BodyMap{}
    gconv.Struct(info, &data)

    err := dao.AlipayMerchantAppConfig.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
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

        // 根据sys_user_id查询商户信息
        employee, err := share_consts.Global.Merchant.Employee().GetEmployeeById(ctx, merchantApp.SysUserId)

        // 3.存储消费者数据并创建用户  kmk-consumer
        getConsumer, err := share_service.Consumer().GetConsumerByAlipayUnionId(ctx, userInfo.Response.UserId)
        consumerRes := &share_model.ConsumerRes{}
        if err != nil && getConsumer == nil {
            shareConsumer := kconv.Struct(userInfo.Response, &share_model.Consumer{})
            shareConsumer.AlipayUnionId = gconv.String(userInfo.Response.UserId)

            if employee != nil {
                shareConsumer.SysUserId = gconv.String(employee.Id)
            }

            consumerRes, err = share_service.Consumer().CreateConsumer(ctx, shareConsumer)
            if err != nil {
                return err
            }
        } else {
            //consumerInfo := share_model.UpdateConsumer{}
            //gconv.Struct(userInfo.Response, &consumerInfo)
            //consumerInfo.Id = merchantApp.SysUserId
            //
            //_, err = share_service.Consumer().UpdateConsumer(ctx, &consumerInfo)
            //if err != nil {
            //	return err
            //}
        }

        // 4.存储阿里消费者记录 alipay-consumer-config
        alipayConsumer, err := service.Consumer().GetConsumerByUserId(ctx, userInfo.Response.UserId)

        if err != nil && alipayConsumer == nil { // 消费者不存在，则创建
            consumerInfo := alipay_model.AlipayConsumerConfig{}
            gconv.Struct(userInfo.Response, &consumerInfo)
            consumerInfo.UserId = gconv.String(userInfo.Response.UserId)

            if employee != nil {
                consumerInfo.SysUserId = employee.Id
            }

            if consumerRes != nil && consumerRes.Id != 0 {
                consumerInfo.SysUserId = consumerRes.Id
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
        platUser, err := share_service.PlatformUser().GetPlatformUserByUserId(ctx, userInfo.Response.UserId)

        if err != nil && platUser == nil { // 不存在创建
            platform := share_model.PlatformUser{
                Id:            idgen.NextId(),
                FacilitatorId: 0,
                OperatorId:    0,
                EmployeeId:    consumerRes.Id,
                MerchantId:    0,
                Platform:      pay_enum.Order.TradeSourceType.Alipay.Code(), // 来源
                ThirdAppId:    merchantApp.ThirdAppId,
                MerchantAppId: merchantApp.AppId,
                UserId:        userInfo.Response.UserId,              // 平台账户唯一标识
                Type:          share_enum.User.Type.Anonymous.Code(), // 用户类型匿名消费者
            }

            if employee != nil {
                platform.EmployeeId = employee.Id
                platform.MerchantId = employee.UnionMainId
            }

            if consumerRes != nil && consumerRes.Id != 0 { // 适用于消费者没有员工的情况下
                platform.EmployeeId = consumerRes.Id
            }

            _, err = share_service.PlatformUser().CreatePlatformUser(ctx, &platform)
            if err != nil {
                return err
            }

        }

        return nil
    })

    if err != nil {
        return false
    }

    return true
}
