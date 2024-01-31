package merchant

import (
	"context"
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
func (s *sUserCertity) AuditConsumer(ctx context.Context, info *alipay_model.CertifyInitReq) (*alipay_model.UserCertifyOpenQueryRes, error) {

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

	ret := alipay_model.UserCertifyOpenInit{}
	gconv.Struct(&res.Response, &ret)

	if ret.Code != gconv.String(10000) {
		re := &alipay_model.UserCertifyOpenQueryRes{}
		re.Response.ErrorResponse = ret.ErrorResponse
		return re, err
	}

	in := make(gopay.BodyMap)
	in.Set("certify_id", ret.CertifyId)

	// 身份认证开始
	_, err = client.UserCertifyOpenCertify(ctx, in)

	if err != nil {
		return nil, err
	}

	// 获取身份认证记录
	re, err := client.UserCertifyOpenQuery(ctx, in)

	r := kconv.Struct(re, &alipay_model.UserCertifyOpenQueryRes{})

	return r, nil

}
