package alipay_hook

import (
	"context"
	"github.com/kysion/alipay-library/alipay_model"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

type ConsumerHookFunc func(ctx context.Context, info interface{}) int64 // 别人订阅我，通常需要返回sys_user_id

type ConsumerHookInfo struct {
	Key   ConsumerKey
	Value ConsumerHookFunc
}
type ConsumerKey struct {
	alipay_enum.Category `json:"category" dc:"业务类别"`
	alipay_enum.ConsumerAction
}
type UserInfo struct {
	SysUserId int64
	alipay_model.UserInfoShare
}
