// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package alipay_service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model"
	enum "github.com/kysion/alipay-library/alipay_model/alipay_enum"
	hook "github.com/kysion/alipay-library/alipay_model/alipay_hook"
	"github.com/kysion/base-library/base_hook"
	"github.com/kysion/gopay/alipay"
	"github.com/kysion/pay-share-library/pay_model"
	"github.com/kysion/pay-share-library/pay_model/pay_enum"
	"github.com/kysion/pay-share-library/pay_model/pay_hook"
)

type (
	ICertify interface {
		// UserCertifyOpenInit 初始化身份认证单据号
		UserCertifyOpenInit(ctx context.Context, appId int64, info *alipay_model.CertifyInitReq) (*alipay_model.UserCertifyOpenInitRes, error)
		// UserCertifyOpenCertify (身份认证开始认证)
		UserCertifyOpenCertify(ctx context.Context, appId int64, certifyId string) (certifyUrl string, err error)
		// UserCertifyOpenQuery 身份认证结果查询
		UserCertifyOpenQuery(ctx context.Context, appId int64, certifyId string) (*alipay_model.UserCertifyOpenQueryRes, error)
	}
	IMerchantTransfer interface {
		// FundTransUniTransfer 单笔转账申请  -
		FundTransUniTransfer(ctx context.Context, appId string, info *alipay_model.FundTransUniTransferReq) (aliRsp *alipay_model.TransUniTransferRes, err error)
		// FundTransCommonQuery 查询转账详情  orderId != OutBizNo  调用转账成功后会返回order_id
		FundTransCommonQuery(ctx context.Context, appId string, outBizNo string) (aliRsp *alipay_model.FundTransCommonQueryRes, err error)
		// FundAccountQuery 余额查询
		FundAccountQuery(ctx context.Context, appId string, userId string) (aliRsp *alipay_model.FundAccountQueryResponse, err error)
	}
	IH5Order interface{}
	IAppAuth interface {
		// AppAuth 具体服务
		AppAuth(ctx context.Context, info g.Map) string
	}
	IAppVersion interface {
		// SubmitVersionAudit 提交应用版本审核
		SubmitVersionAudit(ctx context.Context, info *alipay_model.AppVersionAuditReq, pic *alipay_model.AppVersionAuditPicReq) (*alipay_model.AppVersionAuditRes, error)
		// CancelVersionAudit 撤销版本审核
		CancelVersionAudit(ctx context.Context, version string) (*alipay_model.CancelVersionAuditRes, error)
		// CancelVersion 退回开发版本
		CancelVersion(ctx context.Context, version string) (*alipay_model.CancelVersionRes, error)
		// AppOnline 小程序上架
		AppOnline(ctx context.Context, version string) (*alipay_model.AppOnlineRes, error)
		// AppOffline 小程序下架
		AppOffline(ctx context.Context, version string) (*alipay_model.AppOfflineRes, error)
		// QueryAppVersionList 小程序版本列表查询 https://openapi.alipay.com/v3/alipay/open/mini/version/list/query GET
		QueryAppVersionList(ctx context.Context, versionStatus string) (res *alipay_model.QueryAppVersionListRes, err error)
		// GetAppVersionDetail 小程序版本详情查询 https://openapi.alipay.com/v3/alipay/open/mini/version/detail/query  GET
		GetAppVersionDetail(ctx context.Context, version string) (*alipay_model.QueryAppVersionDetailRes, error)
	}
	IMerchantH5Pay interface {
		// InstallHook 安装Hook的时候，如果状态类型为退款中，需要做响应的退款操作，谨防多模块订阅退款状态，产生重复退款
		InstallHook(actionType pay_enum.OrderStateType, hookFunc pay_hook.OrderHookFunc)
		// TradeCreate H5交易下单 - （当面付）
		TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, merchantApp *alipay_model.AlipayMerchantAppConfig, orderInfo *pay_model.OrderRes, totalAmount float32, userId string) (string, error)
		// H5TradePay H5 支付，返回支付url  -（ 手机网站支付）
		H5TradePay(ctx context.Context, info *alipay_model.TradeOrder, merchantApp *alipay_model.AlipayMerchantAppConfig, orderInfo *pay_model.OrderRes, totalAmount float32) (string, error)
	}
	IMerchantNotify interface {
		// InstallNotifyHook 订阅异步通知Hook
		InstallNotifyHook(hookKey hook.NotifyKey, hookFunc hook.NotifyHookFunc)
		// InstallTradeHook 订阅支付Hook
		InstallTradeHook(hookKey hook.TradeHookKey, hookFunc hook.TradeHookFunc)
		// InstallSubAccountHook 订阅异步通知Hook
		InstallSubAccountHook(hookKey hook.SubAccountHookKey, hookFunc hook.SubAccountHookFunc)
		// MerchantNotifyServices 异步通知地址  用于接收支付宝推送给商户的支付/退款成功的消息。
		MerchantNotifyServices(ctx context.Context) (string, error)
	}
	IMerchantTinyappPay interface {
		// OrderSend 1、发送订单消息   前端完成了
		OrderSend(ctx context.Context)
		// TradeCreate  2、小程序创建交易订单
		TradeCreate(ctx context.Context, info *alipay_model.TradeOrder, merchantApp *alipay_model.AlipayMerchantAppConfig, orderInfo *pay_model.OrderRes, totalAmount float32, userId string) (string, error)
	}
	ISubAccount interface {
		// TradeRelationBind 分账关系绑定
		TradeRelationBind(ctx context.Context, appId int64, info *alipay_model.TradeRelationBindReq) (bool, error)
		// TradeRelationUnbind  分账关系解绑
		TradeRelationUnbind(ctx context.Context, appId string, info *alipay_model.TradeRelationBindReq) (*alipay_model.TradeRelationUnbindResponse, error)
		// TradeRelationBatchQuery 分账关系查询
		TradeRelationBatchQuery(ctx context.Context, appId string, outRequestNo string) (*alipay_model.TradeRoyaltyRelationQueryRes, error)
		// TradeOrderSettleQuery 交易分账查询接口  必须传递第一个参数，或者后两个同时传递
		TradeOrderSettleQuery(ctx context.Context, appId string, settleNo string, outRequestNo string, tradeNo string) (*alipay_model.TradeOrderSettleQueryRes, error)
		// TradeOrderSettle 分账交易下单
		TradeOrderSettle(ctx context.Context, appId string, info alipay_model.TradeOrderSettleReq) (*alipay_model.TradeOrderSettleResponse, error)
	}
	IUserCertity interface {
		// AuditConsumer 身份认证初始化和开始
		AuditConsumer(ctx context.Context, info *alipay_model.CertifyInitReq) (string, error)
	}
	IWallet interface {
		InstallConsumerHook(infoType enum.ConsumerAction, hookFunc hook.ConsumerHookFunc)
		// Wallet 具体服务 H5用户授权 + 小程序
		Wallet(ctx context.Context, info g.Map) string
	}
	IMerchantService interface {
		InstallConsumerHook(infoType hook.ConsumerKey, hookFunc hook.ConsumerHookFunc)
		GetHook() base_hook.BaseHook[hook.ConsumerKey, hook.ConsumerHookFunc]
		// GetUserId 用于检查是否注册,如果已经注册，返会userId
		GetUserId(ctx context.Context, authCode string, appId string) (res string, err error)
		// UserInfoAuth 具体服务 用户授权 + 小程序和H5都兼容
		UserInfoAuth(ctx context.Context, info g.Map) string
	}
	IPayTrade interface {
		// PayTradeCreate  1、创建交易订单   （AppId的H5是没有的，需要写死，小程序有的 ）
		PayTradeCreate(ctx context.Context, info *alipay_model.TradeOrder, userId string, notifyFunc ...hook.NotifyHookFunc) (res string, err error)
		// QueryOrderInfo 查询订单 - 当面付-alipay.trade.query(统一收单线下交易查询)
		QueryOrderInfo(ctx context.Context, outTradeNo string, merchantApp *alipay_model.AlipayMerchantAppConfig) (aliRsp *alipay.TradeQueryResponse, err error)
		// PayTradeClose alipay.trade.close(统一收单交易关闭接口)
		PayTradeClose(ctx context.Context, outTradeNo string, merchantApp *alipay_model.AlipayMerchantAppConfig) (aliRsp *alipay.TradeCloseResponse, err error)
	}
	IUserAuth interface {
		// Cancelled 用户取消授权 （取消关注）
		Cancelled(ctx context.Context, info g.Map) bool
	}
)

var (
	localH5Order            IH5Order
	localAppAuth            IAppAuth
	localAppVersion         IAppVersion
	localMerchantH5Pay      IMerchantH5Pay
	localMerchantNotify     IMerchantNotify
	localMerchantTinyappPay IMerchantTinyappPay
	localSubAccount         ISubAccount
	localUserCertity        IUserCertity
	localWallet             IWallet
	localMerchantService    IMerchantService
	localPayTrade           IPayTrade
	localUserAuth           IUserAuth
	localCertify            ICertify
	localMerchantTransfer   IMerchantTransfer
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

func SubAccount() ISubAccount {
	if localSubAccount == nil {
		panic("implement not found for interface ISubAccount, forgot register?")
	}
	return localSubAccount
}

func RegisterSubAccount(i ISubAccount) {
	localSubAccount = i
}

func UserCertity() IUserCertity {
	if localUserCertity == nil {
		panic("implement not found for interface IUserCertity, forgot register?")
	}
	return localUserCertity
}

func RegisterUserCertity(i IUserCertity) {
	localUserCertity = i
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

func AppVersion() IAppVersion {
	if localAppVersion == nil {
		panic("implement not found for interface IAppVersion, forgot register?")
	}
	return localAppVersion
}

func RegisterAppVersion(i IAppVersion) {
	localAppVersion = i
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

func PayTrade() IPayTrade {
	if localPayTrade == nil {
		panic("implement not found for interface IPayTrade, forgot register?")
	}
	return localPayTrade
}

func RegisterPayTrade(i IPayTrade) {
	localPayTrade = i
}

func UserAuth() IUserAuth {
	if localUserAuth == nil {
		panic("implement not found for interface IUserAuth, forgot register?")
	}
	return localUserAuth
}

func RegisterUserAuth(i IUserAuth) {
	localUserAuth = i
}

func Certify() ICertify {
	if localCertify == nil {
		panic("implement not found for interface ICertify, forgot register?")
	}
	return localCertify
}

func RegisterCertify(i ICertify) {
	localCertify = i
}

func MerchantTransfer() IMerchantTransfer {
	if localMerchantTransfer == nil {
		panic("implement not found for interface IMerchantTransfer, forgot register?")
	}
	return localMerchantTransfer
}

func RegisterMerchantTransfer(i IMerchantTransfer) {
	localMerchantTransfer = i
}

func H5Order() IH5Order {
	if localH5Order == nil {
		panic("implement not found for interface IH5Order, forgot register?")
	}
	return localH5Order
}

func RegisterH5Order(i IH5Order) {
	localH5Order = i
}
