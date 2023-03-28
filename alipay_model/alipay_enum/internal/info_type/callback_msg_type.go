package info_type

import "github.com/kysion/base-library/utility/enum"

type CallBackMsgTypeEnum enum.IEnumCode[string]

// 各种回调消息类型  - 某人某行为产生
type callBackMsgType struct {
	AlipayAppAuth CallBackMsgTypeEnum
	AlipayWallet  CallBackMsgTypeEnum
}

var CallBackMsgType = callBackMsgType{
	AlipayAppAuth: enum.New[CallBackMsgTypeEnum]("alipay_app_auth", "应用认证授权"),
	AlipayWallet:  enum.New[CallBackMsgTypeEnum]("alipay_wallet", "用户登录"),
}
