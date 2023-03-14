package notify

import "github.com/kysion/base-library/utility/enum"

type NotifyTypeEnum enum.IEnumCode[string]

// 各种消息类型
type notifyType struct {
	PayCallBack NotifyTypeEnum
}

var NotifyType = notifyType{
	PayCallBack: enum.New[NotifyTypeEnum]("payCallBack", "支付通知回调"),
}

func (e notifyType) New(code string) NotifyTypeEnum {
	if code == NotifyType.PayCallBack.Code() {
		return NotifyType.PayCallBack
	}
	panic("NotifyTypeEnum: error")
}
