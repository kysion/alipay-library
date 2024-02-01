package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/alipay_utility"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/gopay"
)

type sUserCertity struct{}

func NewUserCertity() service.IUserCertity {

	return &sUserCertity{}

}

func init() {
	service.RegisterUserCertity(NewUserCertity())
}

// AuditConsumer 身份认证初始化和开始
func (s *sUserCertity) AuditConsumer(ctx context.Context, info *alipay_model.CertifyInitReq) (*alipay_model.UserCertifyOpenRes, error) {

	var client *aliyun.AliPay

	// 获取应用ID 根据AppId获取商家相关配置，包括AppAuthToken
	appId := alipay_utility.GetAlipayAppIdFormCtx(ctx)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil || merchantApp == nil {
		return nil, err
	}
	// 判断是否是第三方
	if merchantApp.ThirdAppId != "" {
		client, err = aliyun.NewClient(ctx, appId)
	} else {
		client, err = aliyun.NewMerchantClient(ctx, appId)
	}

	bm := kconv.Struct(&info, &gopay.BodyMap{})

	// 身份认证初始化
	res, err := client.UserCertifyOpenInit(ctx, *bm)

	if err != nil {
		return nil, err
	}

	ret := alipay_model.UserCertifyOpenRes{}
	gconv.Struct(&res.Response, &ret)

	if ret.Code != gconv.String(10000) {
		return nil, err
	}

	in := make(gopay.BodyMap)
	in.Set("certify_id", ret.CertifyId)
	fmt.Println(in)

	// 身份认证开始
	result, err := client.UserCertifyOpenCertify(ctx, in)

	if err != nil {
		return nil, err
	}

	ret.ReturnUrl = result

	return &ret, err
}

// AuditConsumerResponse 查询身份认证结果
func (s *sUserCertity) AuditConsumerResponse(ctx context.Context, certifyId string) (*alipay_model.UserCertifyOpenQueryRes, error) {

	var client *aliyun.AliPay

	// 获取应用ID 根据AppId获取商家相关配置，包括AppAuthToken
	appId := alipay_utility.GetAlipayAppIdFormCtx(ctx)
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil || merchantApp == nil {
		return nil, err
	}
	// 判断是否是第三方
	if merchantApp.ThirdAppId != "" {
		client, err = aliyun.NewClient(ctx, appId)
	} else {
		client, err = aliyun.NewMerchantClient(ctx, appId)
	}

	in := make(gopay.BodyMap)
	in.Set("certify_id", certifyId)

	fmt.Println(in)

	re, err := client.UserCertifyOpenQuery(ctx, in)

	r := kconv.Struct(re, &alipay_model.UserCertifyOpenQueryRes{})

	return r, err

}
