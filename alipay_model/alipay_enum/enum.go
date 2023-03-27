package alipay_enum

import (
	"github.com/kysion/alipay-library/alipay_model/alipay_enum/internal/consumer"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum/internal/info_type"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum/internal/notify"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum/internal/sub_account"
)

type (
	SexType consumer.SexEnum
	// CallbackMsgType 回调消息
	CallbackMsgType info_type.CallBackMsgTypeEnum

	// ServiceNotifyType 应用通知
	ServiceNotifyType info_type.ServiceNotifyTypeEnum

	NotifyType     notify.NotifyTypeEnum
	ConsumerAction consumer.ActionEnum

	// SubAccountBindRes 分账绑定
	SubAccountBindRes sub_account.SubAccountBindResEnum
	// SubAccountAction 分账行为
	SubAccountAction sub_account.OperationTypeEnum
)

var (
	Consumer = consumer.Consumer
	Info     = info_type.Info
	Notify   = notify.Notify

	SubAccount = sub_account.SubAccount
)
