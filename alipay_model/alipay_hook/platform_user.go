package alipay_hook

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/kysion/alipay-library/alipay_model/alipay_enum"
)

type PlatFormUserHookFunc func(ctx context.Context, info PlatformUser) int64

type PlatFormUserHookInfo struct {
	Key   alipay_enum.ConsumerAction
	Value PlatFormUserHookFunc
}

type PlatformUser struct {
	Id            int64       `json:"id"            description:""`
	FacilitatorId int64       `json:"facilitatorId" description:"服务商id"`
	OperatorId    int64       `json:"operatorId"    description:"运营商id"`
	MerchantId    int64       `json:"merchantId"    description:"商户id"`
	EmployeeId    int64       `json:"employeeId"    description:"员工id"`
	Platform      int         `json:"platform"      description:"平台类型：1支付宝、2微信、4抖音、8银联"`
	ThirdAppId    string      `json:"thirdAppId"    description:"第三方平台AppId"`
	MerchantAppId string      `json:"merchantAppId" description:"商家应用AppId"`
	CreatedAt     *gtime.Time `json:"createdAt"     description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     description:""`
	UserId        string      `json:"userId"        description:"平台用户唯一标识"`
	Type          int         `json:"type"          description:"用户类型：0匿名，1用户，2微商，4商户、8广告主、16服务商、32运营中心，64后台"`
}
