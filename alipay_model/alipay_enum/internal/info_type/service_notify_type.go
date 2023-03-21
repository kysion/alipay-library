package info_type

import "github.com/kysion/base-library/utility/enum"

type ServiceNotifyTypeEnum enum.IEnumCode[string]

// 各种消息类型
type serviceNotifyType struct {
	ServiceCheck ServiceNotifyTypeEnum
}

var ServiceNotifyType = serviceNotifyType{
	ServiceCheck: enum.New[ServiceNotifyTypeEnum]("alipay.service.check", "验证应用网关"),
}
