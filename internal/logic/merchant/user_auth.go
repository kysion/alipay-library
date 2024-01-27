package merchant

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/base-library/utility/kconv"
)

type sUserAuth struct {
}

func NewUserAuth() *sUserAuth {
	result := &sUserAuth{}

	result.injectHook()

	return result
}

func (s *sUserAuth) injectHook() {
	//hook := service.Gateway().GetCallbackMsgHook()
	//hook.InstallHook(alipay_enum.Info.CallbackType.AlipayWallet, s.UserInfoAuth)

	// 消息通知
	service.Gateway().GetServiceNotifyTypeHook().InstallHook(alipay_enum.Info.ServiceType.Cancelled, s.Cancelled)
}

// Cancelled 用户取消授权 （取消关注）
func (s *sUserAuth) Cancelled(ctx context.Context, info g.Map) bool {
	/*
		{
		    "version":       "1.1",
		    "msg_method":    "alipay.open.auth.userauth.cancelled",
		    "utc_timestamp": "1706328088729",
		    "sign_type":     "RSA2",
		    "notify_id":     "2024012700262115922056303136304483",
		    "charset":       "UTF-8",
		    "biz_content":   "{\"cancel_time\":1706327962739,\"user_id\":\"2088032632355117\",\"app_id\":\"2021004133631026\"}",
		    "sign":          "MfL3saVl+xjbwr2vp8ObGrCDdQV40i2k6pGFHMHoMZOqloYUikd+OF0unJXHGWel2FomJbXu8hFTs8rPsfEm7/AbcNEKGeiBGQaVcUGOemtWXdkIzihGGtpmBqbq4T4J3WNjwwPAe8B7KXBVFE2ieB5/myKjLu8J9c4wYmaJT9RkShbEPsA7xmOHN4t1hRT012qOM2ypRKRgE9ZNRbhu+cX5jdPRnK9Z9iPPlqtbj3eEhTzb3jwwMIU85BEcIx4bXDI8MPd95s5unvuTQWyILzVPmCO/vs5tXNGFEWDR+BSz/UzgRWANYOtratYhwdb3KXPzkjEpFg3Df38uegG4PA==",
		    "app_id":        "2021004133631026",
		}
	*/
	// 处理用户授权相关
	g.Dump("收到的用户消息事件：", info)
	from := gmap.NewStrAnyMapFrom(info)
	//infoValue := from.Get("info")
	if from.Get("msg_method") != "alipay.open.auth.userauth.cancelled" {
		return true
	}

	bizContent := from.Get("biz_content")
	type Biz struct {
		CancelTime int64  `json:"cancel_time"`
		UserId     string `json:"user_id"`
		AppId      string `json:"app_id"`
	}
	bizData := kconv.Struct(bizContent, &Biz{})

	appId := gconv.String(from.Get("app_id"))
	appConfig, _ := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)

	// 设置用户是否授权
	_, err := service.ConsumerConfig().SetAuthState(ctx, bizData.UserId, appConfig.AppId, alipay_enum.Consumer.AuthState.UnAuth.Code())
	if err != nil {
		return false
	}
	// 输出：success	消息处理成功
	g.RequestFromCtx(ctx).Response.Write("success")
	return true
}
