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

type CertifyResponseReq struct {
	g.Meta    `path:"/:appId/auditConsumerResponse" method:"post" summary:"实名认证结果" tags:"Alipay消费者应用"`
	CertifyId string `json:"certify_id,omitempty"`
}
