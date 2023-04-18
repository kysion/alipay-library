// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
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
		GetConsumerById(ctx context.Context, id int64) (*alipay_model.AlipayConsumerConfig, error)
		GetConsumerByUserId(ctx context.Context, userId string) (*alipay_model.AlipayConsumerConfig, error)
		GetConsumerBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayConsumerConfig, error)
		CreateConsumer(ctx context.Context, info *alipay_model.AlipayConsumerConfig) (*alipay_model.AlipayConsumerConfig, error)
		UpdateConsumer(ctx context.Context, id int64, info *alipay_model.UpdateConsumerReq) (bool, error)
		UpdateConsumerState(ctx context.Context, id int64, state int) (bool, error)
	}
	IGateway interface {
		GetCallbackMsgHook() *base_hook.BaseHook[enum.CallbackMsgType, hook.ServiceMsgHookFunc]
		GetServiceNotifyTypeHook() base_hook.BaseHook[enum.ServiceNotifyType, hook.ServiceNotifyHookFunc]
		GatewayServices(ctx context.Context) (string, error)
		GatewayCallback(ctx context.Context) (string, error)
	}
	IMerchantAppConfig interface {
		GetMerchantAppConfigById(ctx context.Context, id int64) (*alipay_model.AlipayMerchantAppConfig, error)
		GetMerchantAppConfigByAppId(ctx context.Context, id string) (*alipay_model.AlipayMerchantAppConfig, error)
		GetMerchantAppConfigBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayMerchantAppConfig, error)
		CreateMerchantAppConfig(ctx context.Context, info *alipay_model.AlipayMerchantAppConfig) (*alipay_model.AlipayMerchantAppConfig, error)
		UpdateMerchantAppConfig(ctx context.Context, info *alipay_model.UpdateMerchantAppConfigReq) (bool, error)
		UpdateState(ctx context.Context, id int64, state int) (bool, error)
		UpdateAppAuthToken(ctx context.Context, info *alipay_model.UpdateMerchantAppAuthToken) (bool, error)
		UpdateAppConfigHttps(ctx context.Context, info *alipay_model.UpdateMerchantAppConfigHttpsReq) (bool, error)
		UpdateMerchantKeyCert(ctx context.Context, info *alipay_model.UpdateMerchantKeyCertReq) (bool, error)
		CreatePolicy(ctx context.Context, info *alipay_model.CreatePolicyReq) (bool, error)
		GetPolicy(ctx context.Context, appId string) (*alipay_model.GetPolicyRes, error)
	}
	IThirdAppConfig interface {
		GetThirdAppConfigById(ctx context.Context, id int64) (*alipay_model.AlipayThirdAppConfig, error)
		GetThirdAppConfigByAppId(ctx context.Context, id string) (*alipay_model.AlipayThirdAppConfig, error)
		GetThirdAppConfigBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayThirdAppConfig, error)
		CreateThirdAppConfig(ctx context.Context, info *alipay_model.AlipayThirdAppConfig) (*alipay_model.AlipayThirdAppConfig, error)
		UpdateThirdAppConfig(ctx context.Context, info *alipay_model.UpdateThirdAppConfig) (bool, error)
		UpdateState(ctx context.Context, id int64, state int) (bool, error)
		UpdateAppAuthToken(ctx context.Context, info *alipay_model.UpdateThirdAppAuthToken) (bool, error)
		UpdateAppConfigHttps(ctx context.Context, info *alipay_model.UpdateThirdAppConfigHttpsReq) (bool, error)
		UpdateThirdKeyCert(ctx context.Context, info *alipay_model.UpdateThirdKeyCertReq) (bool, error)
	}
)

var (
	localMerchantAppConfig IMerchantAppConfig
	localThirdAppConfig    IThirdAppConfig
	localConsumerConfig    IConsumerConfig
	localGateway           IGateway
)

func ThirdAppConfig() IThirdAppConfig {
	if localThirdAppConfig == nil {
		panic("implement not found for interface IThirdAppConfig, forgot register?")
	}
	return localThirdAppConfig
}

func RegisterThirdAppConfig(i IThirdAppConfig) {
	localThirdAppConfig = i
}

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
