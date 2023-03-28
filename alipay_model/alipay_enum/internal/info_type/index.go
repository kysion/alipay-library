package info_type

type info struct {
	CallbackType callBackMsgType

	ServiceType serviceNotifyType
}

var Info = info{
	CallbackType: CallBackMsgType,
	ServiceType:  ServiceNotifyType,
}
