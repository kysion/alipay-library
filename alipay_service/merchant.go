// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package alipay_service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_model/pay_hook"
)

type (
	IMerchantService interface {
		InstallConsumerHook(infoType alipay_enum.ConsumerAction, hookFunc hook.ConsumerHookFunc)
		GetHook() base_hook.BaseHook[alipay_enum.ConsumerAction, hook.ConsumerHookFunc]
		GetUserId(ctx context.Context, authCode string, appId string) (res string, err error)
		UserInfoAuth(ctx context.Context, authCode string, appId string, sysUserId ...int64) (res *alipay_model.UserInfoShare, err error)
	}
	IMerchantTinyappPay interface {
		OrderSend(ctx context.Context)
		TradeCreate(ctx context.Context, info *alipay_model.CreateTrade) (aliRsp *alipay_model.TradeCreateResponse, err error)
	}
	IWallet interface {
		InstallConsumerHook(infoType enum.ConsumerAction, hookFunc hook.ConsumerHookFunc)
		Wallet(ctx context.Context, info g.Map) bool
	}
	IAppAuth interface {
		AppAuth(ctx context.Context, info g.Map) bool
	}
	IMerchantH5Pay interface {
		InstallHook(actionType pay_enum.OrderStateType, hookFunc pay_hook.OrderHookFunc)
		H5TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, notifyFunc ...hook.NotifyHookFunc)
		QueryOrderInfo(ctx context.Context, outTradeNo string, merchantAppId string, thirdAppId string, appAuthToken string)
	}
	IMerchantNotify interface {
		InstallNotifyHook(hookKey hook.NotifyKey, hookFunc hook.NotifyHookFunc)
		InstallTradeHook(hookKey hook.TradeHookKey, hookFunc hook.TradeHookFunc)
		MerchantNotifyServices(ctx context.Context) (string, error)
	}
)

var (
	localWallet             IWallet
	localAppAuth            IAppAuth
	localMerchantH5Pay      IMerchantH5Pay
	localMerchantNotify     IMerchantNotify
	localMerchantService    IMerchantService
	localMerchantTinyappPay IMerchantTinyappPay
)

func MerchantTinyappPay() IMerchantTinyappPay {
	if localMerchantTinyappPay == nil {
		panic("implement not found for interface IMerchantTinyappPay, forgot register?")
	}
	return localMerchantTinyappPay
}

func RegisterMerchantTinyappPay(i IMerchantTinyappPay) {
	localMerchantTinyappPay = i
}

func Wallet() IWallet {
	if localWallet == nil {
		panic("implement not found for interface IWallet, forgot register?")
	}
	return localWallet
}

func RegisterWallet(i IWallet) {
	localWallet = i
}

func AppAuth() IAppAuth {
	if localAppAuth == nil {
		panic("implement not found for interface IAppAuth, forgot register?")
	}
	return localAppAuth
}

func RegisterAppAuth(i IAppAuth) {
	localAppAuth = i
}

func MerchantH5Pay() IMerchantH5Pay {
	if localMerchantH5Pay == nil {
		panic("implement not found for interface IMerchantH5Pay, forgot register?")
	}
	return localMerchantH5Pay
}

func RegisterMerchantH5Pay(i IMerchantH5Pay) {
	localMerchantH5Pay = i
}

func MerchantNotify() IMerchantNotify {
	if localMerchantNotify == nil {
		panic("implement not found for interface IMerchantNotify, forgot register?")
	}
	return localMerchantNotify
}

func RegisterMerchantNotify(i IMerchantNotify) {
	localMerchantNotify = i
}

func MerchantService() IMerchantService {
	if localMerchantService == nil {
		panic("implement not found for interface IMerchantService, forgot register?")
	}
	return localMerchantService
}

func RegisterMerchantService(i IMerchantService) {
	localMerchantService = i
}
