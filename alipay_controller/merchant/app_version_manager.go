package merchant_controller

import (
	"context"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/api/alipay_v1/alipay_merchant_v1"
)

var AppVersionManager = cAppVersionManager{}

type cAppVersionManager struct{}

// SubmitAppVersionAudit 提交应用版本审核
func (c *cAppVersionManager) SubmitAppVersionAudit(ctx context.Context, req *alipay_merchant_v1.SubmitAppVersionAuditReq) (*alipay_model.AppVersionAuditRes, error) {
	ret, err := service.AppVersion().SubmitVersionAudit(ctx, &req.AppVersionAuditReq, &req.AppVersionAuditPicReq)

	return ret, err
}

// CancelAppVersionAudit 撤销应用版本审核
func (c *cAppVersionManager) CancelAppVersionAudit(ctx context.Context, req *alipay_merchant_v1.CancelAppVersionAuditReq) (*alipay_model.CancelVersionAuditRes, error) {
	ret, err := service.AppVersion().CancelVersionAudit(ctx, req.AppVersion)

	return ret, err
}

// CancelAppVersion 退回开发版本
func (c *cAppVersionManager) CancelAppVersion(ctx context.Context, req *alipay_merchant_v1.CancelAppVersionReq) (*alipay_model.CancelVersionRes, error) {
	ret, err := service.AppVersion().CancelVersion(ctx, req.AppVersion)

	return ret, err
}

// AppOnline 小程序上架
func (c *cAppVersionManager) AppOnline(ctx context.Context, req *alipay_merchant_v1.AppOnlineReq) (*alipay_model.AppOnlineRes, error) {
	ret, err := service.AppVersion().AppOnline(ctx, req.AppVersion)

	return ret, err
}

// AppOffline 小程序下架
func (c *cAppVersionManager) AppOffline(ctx context.Context, req *alipay_merchant_v1.AppOfflineReq) (*alipay_model.AppOfflineRes, error) {
	ret, err := service.AppVersion().AppOffline(ctx, req.AppVersion)

	return ret, err
}

// QueryAppVersionList 查询小程序版本列表
func (c *cAppVersionManager) QueryAppVersionList(ctx context.Context, req *alipay_merchant_v1.QueryAppVersionListReq) (*alipay_model.QueryAppVersionListRes, error) {
	ret, err := service.AppVersion().QueryAppVersionList(ctx, req.VersionStatus)

	return ret, err
}

// GetAppVersionDetail 查询小程序版本详情
func (c *cAppVersionManager) GetAppVersionDetail(ctx context.Context, req *alipay_merchant_v1.GetAppVersionDetailReq) (*alipay_model.QueryAppVersionDetailRes, error) {
	ret, err := service.AppVersion().GetAppVersionDetail(ctx, req.AppVersion)

	return ret, err
}
