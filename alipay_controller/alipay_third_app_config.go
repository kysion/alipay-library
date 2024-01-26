package alipay_controller

import (
    "context"
    "github.com/SupenBysz/gf-admin-community/api_v1"
    "github.com/SupenBysz/gf-admin-community/utility/funs"
    "github.com/kysion/alipay-library/alipay_api/alipay_v1/alipay_third_v1"
    "github.com/kysion/alipay-library/alipay_service"
)

var AlipayThirdAppConfig = cAlipayThirdAppConfig{}

type cAlipayThirdAppConfig struct{}

// UpdateState 修改状态
func (s *cAlipayThirdAppConfig) UpdateState(ctx context.Context, req *alipay_third_v1.UpdateStateReq) (api_v1.BoolRes, error) {
	return funs.CheckPermission(ctx,
		func() (api_v1.BoolRes, error) {
			ret, err := alipay_service.ThirdAppConfig().UpdateState(ctx, req.Id, req.State)
			return ret == true, err
		},
		// 记得添加权限
		// alipay_permission.ThirdAppConfig.PermissionType.Update,
	)
}

// CreateThirdAppConfig  创建第三方应用配置信息
func (s *cAlipayThirdAppConfig) CreateThirdAppConfig(ctx context.Context, req *alipay_third_v1.CreateThirdAppConfigReq) (*alipay_third_v1.ThirdAppConfigRes, error) {
	//return funs.CheckPermission(ctx,
	//	func() (*v1.ThirdAppConfigRes, error) {
	//		ret, err := alipay_service.ThirdAppConfig().CreateThirdAppConfig(ctx, &req.AlipayThirdAppConfig)
	//		return (*v1.ThirdAppConfigRes)(ret), err
	//	},
	//	// 记得添加权限
	//	// alipay_permission.ThirdAppConfig.PermissionType.Update,
	//)

	ret, err := alipay_service.ThirdAppConfig().CreateThirdAppConfig(ctx, &req.AlipayThirdAppConfig)
	return (*alipay_third_v1.ThirdAppConfigRes)(ret), err
}

// GetThirdAppConfigByAppId 根据AppId查找第三方应用配置信息
func (s *cAlipayThirdAppConfig) GetThirdAppConfigByAppId(ctx context.Context, req *alipay_third_v1.GetThirdAppConfigByIdReq) (*alipay_third_v1.ThirdAppConfigRes, error) {
	return funs.CheckPermission(ctx,
		func() (*alipay_third_v1.ThirdAppConfigRes, error) {
			ret, err := alipay_service.ThirdAppConfig().GetThirdAppConfigById(ctx, req.Id)
			return (*alipay_third_v1.ThirdAppConfigRes)(ret), err
		},
		// 记得添加权限
		// alipay_permission.ThirdAppConfig.PermissionType.Update,
	)

}

// UpdateAppConfig 修改服务商基础信息
func (s *cAlipayThirdAppConfig) UpdateAppConfig(ctx context.Context, req *alipay_third_v1.UpdateThirdAppConfigReq) (api_v1.BoolRes, error) {
	return funs.CheckPermission(ctx,
		func() (api_v1.BoolRes, error) {
			ret, err := alipay_service.ThirdAppConfig().UpdateThirdAppConfig(ctx, &req.UpdateThirdAppConfig)
			return ret == true, err
		},
		// 记得添加权限
		// alipay_permission.ThirdAppConfig.PermissionType.Update,
	)
	//
	//ret, err := alipay_service.ThirdAppConfig().UpdateAppConfig(ctx, &req.UpdateThirdAppConfigReq)
	//
	//return ret == true, err
}

// UpdateThirdAppConfigHttps 修改服务商应用Https配置
func (s *cAlipayThirdAppConfig) UpdateThirdAppConfigHttps(ctx context.Context, req *alipay_third_v1.UpdateThirdAppConfigHttpsReq) (api_v1.BoolRes, error) {
	ret, err := alipay_service.ThirdAppConfig().UpdateAppConfigHttps(ctx, &req.UpdateThirdAppConfigHttpsReq)
	return ret == true, err
}

// UpdateThirdKeyCertConfig 修改密钥证书
func (c *cAlipayThirdAppConfig) UpdateThirdKeyCertConfig(ctx context.Context, req *alipay_third_v1.UpdateThirdKeyCertReq) (api_v1.BoolRes, error) {
	ret, err := alipay_service.ThirdAppConfig().UpdateThirdKeyCert(ctx, &req.UpdateThirdKeyCertReq)
	return ret == true, err
}
