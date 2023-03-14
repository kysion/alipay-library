package alipay_enum

import (
	"github.com/kysion/alipay-test/alipay_model/alipay_enum/internal/consumer"
	"github.com/kysion/alipay-test/alipay_model/alipay_enum/internal/info_type"
	"github.com/kysion/alipay-test/alipay_model/alipay_enum/internal/notify"
)

type (
	SexType    = consumer.SexEnum
	InfoType   info_type.InfoTypeEnum
	NotifyType notify.NotifyTypeEnum
)

var (
	Consumer = consumer.Consumer
	Info     = info_type.Info
	Notify   = notify.Notify
)
