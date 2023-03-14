package alipay_controller

import (
    "context"
    "github.com/SupenBysz/gf-admin-community/api_v1"
    "github.com/SupenBysz/gf-admin-community/utility/funs"
    "github.com/kysion/alipay-test/alipay_service"
    alipay_merchant_v1 "github.com/kysion/alipay-test/api/alipay_v1/alipay_merchant_v1"
)

var AlipayMerchantAppConfig = cAlipayMerchantAppConfig{}

type cAlipayMerchantAppConfig struct{}

// UpdateState 修改状态
func (s *cAlipayMerchantAppConfig) UpdateState(ctx context.Context, req *alipay_merchant_v1.UpdateStateReq) (api_v1.BoolRes, error) {
    return funs.CheckPermission(ctx,
        func() (api_v1.BoolRes, error) {
            ret, err := alipay_service.MerchantAppConfig().UpdateState(ctx, req.Id, req.State)
            return ret == true, err
        },
        // 记得添加权限
        // alipay_permission.MerchantAppConfig.PermissionType.Update,
    )
}

// CreateMerchantAppConfig  创建第三方应用配置信息
func (s *cAlipayMerchantAppConfig) CreateMerchantAppConfig(ctx context.Context, req *alipay_merchant_v1.CreateMerchantAppConfigReq) (*alipay_merchant_v1.MerchantAppConfigRes, error) {
    //return funs.CheckPermission(ctx,
    //	func() (*v1.MerchantAppConfigRes, error) {
    //		ret, err := alipay_service.MerchantAppConfig().CreateMerchantAppConfig(ctx, &req.AlipayMerchantAppConfig)
    //		return (*v1.MerchantAppConfigRes)(ret), err
    //	},
    //	// 记得添加权限
    //	// alipay_permission.MerchantAppConfig.PermissionType.Update,
    //)

    ret, err := alipay_service.MerchantAppConfig().CreateMerchantAppConfig(ctx, &req.AlipayMerchantAppConfig)
    return (*alipay_merchant_v1.MerchantAppConfigRes)(ret), err
}

// GetMerchantAppConfigByAppId 根据AppId查找第三方应用配置信息
func (s *cAlipayMerchantAppConfig) GetMerchantAppConfigByAppId(ctx context.Context, req *alipay_merchant_v1.GetMerchantAppConfigByIdReq) (*alipay_merchant_v1.MerchantAppConfigRes, error) {
    return funs.CheckPermission(ctx,
        func() (*alipay_merchant_v1.MerchantAppConfigRes, error) {
            ret, err := alipay_service.MerchantAppConfig().GetMerchantAppConfigById(ctx, req.Id)
            return (*alipay_merchant_v1.MerchantAppConfigRes)(ret), err
        },
        // 记得添加权限
        // alipay_permission.MerchantAppConfig.PermissionType.Update,
    )

}

// UpdateAppConfig 修改服务商基础信息
func (s *cAlipayMerchantAppConfig) UpdateAppConfig(ctx context.Context, req *alipay_merchant_v1.UpdateMerchantAppConfigReq) (api_v1.BoolRes, error) {
    return funs.CheckPermission(ctx,
        func() (api_v1.BoolRes, error) {
            ret, err := alipay_service.MerchantAppConfig().UpdateMerchantAppConfig(ctx, &req.UpdateMerchantAppConfigReq)
            return ret == true, err
        },
        // 记得添加权限
        // alipay_permission.MerchantAppConfig.PermissionType.Update,
    )
    //
    //ret, err := alipay_service.MerchantAppConfig().UpdateAppConfig(ctx, &req.UpdateMerchantAppConfigReq)
    //
    //return ret == true, err
}

// UpdateMerchantAppConfigHttps 修改服务商应用Https配置
func (s *cAlipayMerchantAppConfig) UpdateMerchantAppConfigHttps(ctx context.Context, req *alipay_merchant_v1.UpdateMerchantAppConfigHttpsReq) (api_v1.BoolRes, error) {
    ret, err := alipay_service.MerchantAppConfig().UpdateAppConfigHttps(ctx, &req.UpdateMerchantAppConfigHttpsReq)
    return ret == true, err
}

// UpdateMerchantKeyCertConfig 修改密钥证书
func (c *cAlipayMerchantAppConfig) UpdateMerchantKeyCertConfig(ctx context.Context, req *alipay_merchant_v1.UpdateMerchantKeyCertReq) (api_v1.BoolRes, error) {
    ret, err := alipay_service.MerchantAppConfig().UpdateMerchantKeyCert(ctx, &req.UpdateMerchantKeyCertReq)
    return ret == true, err
}
