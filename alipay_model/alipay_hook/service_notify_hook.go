package alipay_hook

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

// 应用消息 - 平台主动发的

type ServiceNotifyHookFunc func(ctx context.Context, info g.Map) bool

type ServiceNotifyHookInfo struct {
	Key   alipay_enum.ServiceNotifyType
	Value ServiceNotifyHookFunc
}
