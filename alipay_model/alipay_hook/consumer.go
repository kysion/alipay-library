package alipay_hook

import (
	"context"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

type ConsumerHookFunc func(ctx context.Context, info interface{}) int64

type ConsumerHookInfo struct {
	Key   alipay_enum.ConsumerAction
	Value ConsumerHookFunc
}

type UserInfo struct {
	SysUserId int64
	alipay_model.UserInfoShare
}
