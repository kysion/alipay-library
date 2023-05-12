package alipay_hook

import (
	"context"
	entity "github.com/kysion/alipay-library/alipay_model/alipay_entity"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

type PlatFormUserHookFunc func(ctx context.Context, info entity.PlatformUser) int64

type PlatFormUserHookInfo struct {
	Key   alipay_enum.ConsumerAction
	Value PlatFormUserHookFunc
}
