package alipay_hook

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

// 应用通知 - 由某人产生
type ServiceMsgHookFunc func(ctx context.Context, info g.Map) bool

type ServiceMsgHookInfo struct {
	Key   alipay_enum.CallbackMsgType
	Value ServiceMsgHookFunc
}