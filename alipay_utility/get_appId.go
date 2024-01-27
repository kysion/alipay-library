package alipay_utility

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
)

func GetAlipayAppIdFormCtx(ctx context.Context) string {
	// appId的32进制编码
	appId := g.RequestFromCtx(ctx).Get("appId").String()

	if len(appId) < 16 {
		id, _ := strconv.ParseInt(appId, 36, 0)

		appId = gconv.String(id)
	}

	return appId
}
