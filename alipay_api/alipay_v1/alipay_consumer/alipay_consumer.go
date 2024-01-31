package alipay_consumer_v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model"
)

// 阿里消费者相关接口

type CertifyReq struct {
	g.Meta `path:"/:appId/auditConsumer" method:"post" summary:"实名认证" tags:"Alipay消费者应用"`
	alipay_model.CertifyInitReq
}
