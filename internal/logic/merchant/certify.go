package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/gopay"
)

// 实名认证相关
type sCertify struct{}

func NewCertify() *sCertify {
	return &sCertify{}
}

// UserCertifyOpenInit 初始化身份认证单据号
func (s *sCertify) UserCertifyOpenInit(ctx context.Context, appId int64, info *alipay_model.CertifyInitReq) (*alipay_model.UserCertifyOpenInitRes, error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay实名认证初始化 ------- ", "Alipay-sCertify")

	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(appId))
	if err != nil {
		return nil, err
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)

	data := kconv.Struct(info, &gopay.BodyMap{})
	data.Set("app_auth_token", merchantApp.AppAuthToken)

	// 初始化身份认证
	initRes, err := client.UserCertifyOpenInit(ctx, *data)
	if err != nil {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "初始化身份认证单据号失败", "Alipay-sCertify")
	}

	res := kconv.Struct(initRes, &alipay_model.UserCertifyOpenInitRes{})

	return res, nil
}

// UserCertifyOpenCertify (身份认证开始认证)
func (s *sCertify) UserCertifyOpenCertify(ctx context.Context, appId int64, certifyId string) (certifyUrl string, err error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay实名认证开始 ------- ", "Alipay-sCertify")

	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(appId))
	if err != nil {
		return "", err
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)

	// 开始身份认证
	openRes, err := client.UserCertifyOpenCertify(ctx, gopay.BodyMap{
		"certify_id": certifyId,
	})

	if err != nil {
		return "", sys_service.SysLogs().ErrorSimple(ctx, err, "启动身份认证失败", "Alipay-sCertify")
	}

	return openRes, nil
}

// UserCertifyOpenQuery 身份认证结果查询
func (s *sCertify) UserCertifyOpenQuery(ctx context.Context, appId int64, certifyId string) (*alipay_model.UserCertifyOpenQueryRes, error) {
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------Alipay实名认证初始化 ------- ", "Alipay-sCertify")

	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(appId))
	if err != nil {
		return nil, err
	}

	// 通过商家中的第三方应用的AppId创建客户端
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)

	// 身份认证结果查询
	queryRes, err := client.UserCertifyOpenQuery(ctx, gopay.BodyMap{
		"certify_id": certifyId,
	})

	if err != nil {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "初始化身份认证单据号失败", "Alipay-sCertify")
	}

	res := kconv.Struct(queryRes, &alipay_model.UserCertifyOpenQueryRes{})

	return res, nil
}
