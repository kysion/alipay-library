package hook

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-test/alipay_model/alipay_enum"
)

type ServiceMsgHookFunc func(ctx context.Context, info g.Map) bool

type ServiceMsgHookInfo struct {
	Key   alipay_enum.InfoType
	Value ServiceMsgHookFunc
}
