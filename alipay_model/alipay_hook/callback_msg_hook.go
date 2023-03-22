package alipay_hook

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

// ServiceMsgHookFunc 应用消息 - 由某人产生  对应回调CallBack
type ServiceMsgHookFunc func(ctx context.Context, info g.Map) string // 通常需要返回用户userId

type ServiceMsgHookInfo struct {
	Key   alipay_enum.CallbackMsgType
	Value ServiceMsgHookFunc
}
