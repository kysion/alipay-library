package info_type

import "github.com/kysion/base-library/utility/enum"

type InfoTypeEnum enum.IEnumCode[string]

// 各种消息类型
type infoType struct {
	Ticket                InfoTypeEnum
	AuthorizerAccessToken InfoTypeEnum
	AlipayAppAuth         InfoTypeEnum
	AlipayWallet          InfoTypeEnum
}

var InfoType = infoType{
	Ticket:                enum.New[InfoTypeEnum]("Ticket", "票据"),
	AuthorizerAccessToken: enum.New[InfoTypeEnum]("AuthorizerAccessToken", "授权小程序接口调用凭据"),
	AlipayAppAuth:         enum.New[InfoTypeEnum]("alipay_app_auth", "应用认证授权"),
	AlipayWallet:          enum.New[InfoTypeEnum]("alipay_wallet", "用户登录"),
}
