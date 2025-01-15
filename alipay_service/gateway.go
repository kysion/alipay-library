// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package alipay_service

import (
	"context"

	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	"github.com/kysion/base-library/base_hook"
)

type (
	IConsumerConfig interface {
		// GetConsumerById 根据id查找消费者信息
		GetConsumerById(ctx context.Context, id int64) (*alipay_model.AlipayConsumerConfig, error)
		// GetConsumerByUserId  根据平台用户id查询消费者信息
		GetConsumerByUserId(ctx context.Context, userId string) (*alipay_model.AlipayConsumerConfig, error)
		// GetConsumerByUserIdAndAppId  根据平台用户id+ AppId查询消费者信息
		GetConsumerByUserIdAndAppId(ctx context.Context, userId string, appId string) (*alipay_model.AlipayConsumerConfig, error)
		// GetConsumerBySysUserId  根据用户id查询消费者信息
		GetConsumerBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayConsumerConfig, error)
		// CreateConsumer  创建消费者信息
		CreateConsumer(ctx context.Context, info *alipay_model.AlipayConsumerConfig) (*alipay_model.AlipayConsumerConfig, error)
		// UpdateConsumer 更新消费者信息
		UpdateConsumer(ctx context.Context, id int64, info *alipay_model.UpdateConsumerReq) (bool, error)
		// UpdateConsumerState 修改用户状态
		UpdateConsumerState(ctx context.Context, id int64, state int) (bool, error)
		// SetAuthState 是否授权
		SetAuthState(ctx context.Context, userId string, appID string, authState int) (bool, error)
	}
	IGateway interface {
		// GetCallbackMsgHook 返回回调消息Hook对象
		GetCallbackMsgHook() *base_hook.BaseHook[enum.CallbackMsgType, hook.ServiceMsgHookFunc]
		GetServiceNotifyTypeHook() *base_hook.BaseHook[enum.ServiceNotifyType, hook.ServiceNotifyHookFunc]
		// GatewayServices 接收消息通知  B端消息
		GatewayServices(ctx context.Context) (string, error)
		// GatewayCallback 接收消息回调  C端消息
		GatewayCallback(ctx context.Context) (string, error)
	}
	IMerchantAppConfig interface {
		// GetMerchantAppConfigById 根据id查找商家应用配置信息
		GetMerchantAppConfigById(ctx context.Context, id int64) (*alipay_model.AlipayMerchantAppConfig, error)
		// GetMerchantAppConfigByAppId 根据AppId查找商家应用配置信息
		GetMerchantAppConfigByAppId(ctx context.Context, id string) (*alipay_model.AlipayMerchantAppConfig, error)
		// GetMerchantAppConfigBySysUserId  根据商家id查询商家应用配置信息
		GetMerchantAppConfigBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayMerchantAppConfig, error)
		// CreateMerchantAppConfig  创建商家应用配置信息
		CreateMerchantAppConfig(ctx context.Context, info *alipay_model.AlipayMerchantAppConfig) (*alipay_model.AlipayMerchantAppConfig, error)
		// UpdateMerchantAppConfig 更新商家应用配置信息
		UpdateMerchantAppConfig(ctx context.Context, info *alipay_model.UpdateMerchantAppConfigReq) (bool, error)
		// UpdateState 修改商家应用状态
		UpdateState(ctx context.Context, id int64, state int) (bool, error)
		// UpdateAppAuthToken 更新Token  商家应用授权token
		UpdateAppAuthToken(ctx context.Context, info *alipay_model.UpdateMerchantAppAuthToken) (bool, error)
		// UpdateAppConfigHttps 修改商家应用Https配置
		UpdateAppConfigHttps(ctx context.Context, info *alipay_model.UpdateMerchantAppConfigHttpsReq) (bool, error)
		// UpdateMerchantKeyCert 更新商家应用配置证书密钥
		UpdateMerchantKeyCert(ctx context.Context, info *alipay_model.UpdateMerchantKeyCertReq) (bool, error)
		// CreatePolicy 创建用户协议或隐私协议
		CreatePolicy(ctx context.Context, info *alipay_model.CreatePolicyReq) (bool, error)
		// GetPolicy 获取协议
		GetPolicy(ctx context.Context, appId string) (*alipay_model.GetPolicyRes, error)
	}
	IThirdAppConfig interface {
		// GetThirdAppConfigById 根据id查找第三方应用配置信息
		GetThirdAppConfigById(ctx context.Context, id int64) (*alipay_model.AlipayThirdAppConfig, error)
		// GetThirdAppConfigByAppId 根据AppId查找第三方应用配置信息
		GetThirdAppConfigByAppId(ctx context.Context, id string) (*alipay_model.AlipayThirdAppConfig, error)
		// GetThirdAppConfigBySysUserId  根据用户id查询第三方应用配置信息
		GetThirdAppConfigBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayThirdAppConfig, error)
		// CreateThirdAppConfig  创建第三方应用配置信息
		CreateThirdAppConfig(ctx context.Context, info *alipay_model.AlipayThirdAppConfig) (*alipay_model.AlipayThirdAppConfig, error)
		// UpdateThirdAppConfig 更新第三方应用基础配置信息
		UpdateThirdAppConfig(ctx context.Context, info *alipay_model.UpdateThirdAppConfig) (bool, error)
		// UpdateState 修改第三方应用状态
		UpdateState(ctx context.Context, id int64, state int) (bool, error)
		// UpdateAppAuthToken 更新Token  服务商应用授权token
		UpdateAppAuthToken(ctx context.Context, info *alipay_model.UpdateThirdAppAuthToken) (bool, error)
		// UpdateAppConfigHttps 修改服务商应用Https配置
		UpdateAppConfigHttps(ctx context.Context, info *alipay_model.UpdateThirdAppConfigHttpsReq) (bool, error)
		// UpdateThirdKeyCert 更新第三方应用配置证书密钥
		UpdateThirdKeyCert(ctx context.Context, info *alipay_model.UpdateThirdKeyCertReq) (bool, error)
	}
)

var (
	localConsumerConfig    IConsumerConfig
	localGateway           IGateway
	localMerchantAppConfig IMerchantAppConfig
	localThirdAppConfig    IThirdAppConfig
)

func ConsumerConfig() IConsumerConfig {
	if localConsumerConfig == nil {
		panic("implement not found for interface IConsumerConfig, forgot register?")
	}
	return localConsumerConfig
}

func RegisterConsumerConfig(i IConsumerConfig) {
	localConsumerConfig = i
}

func Gateway() IGateway {
	if localGateway == nil {
		panic("implement not found for interface IGateway, forgot register?")
	}
	return localGateway
}

func RegisterGateway(i IGateway) {
	localGateway = i
}

func MerchantAppConfig() IMerchantAppConfig {
	if localMerchantAppConfig == nil {
		panic("implement not found for interface IMerchantAppConfig, forgot register?")
	}
	return localMerchantAppConfig
}

func RegisterMerchantAppConfig(i IMerchantAppConfig) {
	localMerchantAppConfig = i
}

func ThirdAppConfig() IThirdAppConfig {
	if localThirdAppConfig == nil {
		panic("implement not found for interface IThirdAppConfig, forgot register?")
	}
	return localThirdAppConfig
}

func RegisterThirdAppConfig(i IThirdAppConfig) {
	localThirdAppConfig = i
}
