package alipay_merchant_v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/alipay-library/alipay_model"
)

type SubmitAppVersionAuditReq struct {
	g.Meta `path:"/:appId/submitAppVersionAudit" method:"get" summary:"提交应用版本审核" tags:"Alipay小程序管理"`
	SubmitVersion
}

type SubmitVersion struct {
	alipay_model.AppVersionAuditReq
	alipay_model.AppVersionAuditPicReq
}

type CancelAppVersionAuditReq struct {
	g.Meta     `path:"/:appId/cancelAppVersionAudit" method:"get" summary:"撤销应用版本审核" tags:"Alipay小程序管理"`
	AppVersion string `json:"app_version" dc:"版本号"`
}

type CancelAppVersionReq struct {
	g.Meta     `path:"/:appId/cancelAppVersion" method:"get" summary:"退回开发版本" tags:"Alipay小程序管理"`
	AppVersion string `json:"app_version" dc:"版本号"`
}

type AppOnlineReq struct {
	g.Meta     `path:"/:appId/appOnline" method:"get" summary:"小程序上架" tags:"Alipay小程序管理"`
	AppVersion string `json:"app_version" dc:"版本号"`
}

type AppOfflineReq struct {
	g.Meta     `path:"/:appId/appOffline" method:"get" summary:"小程序下架" tags:"Alipay小程序管理"`
	AppVersion string `json:"app_version" dc:"版本号"`
}

type QueryAppVersionListReq struct {
	g.Meta `path:"/:appId/queryAppVersionList" method:"post" summary:"小程序版本列表查询" tags:"Alipay小程序管理"`
	// query请求URL传参参数
	//BundleId      string `json:"bundle_id" dc:"端参数"`
	VersionStatus string `json:"version_status" dc:"版本状态列表，用英文逗号,分割，不填默认不返回，说明如下：INIT: 开发中, AUDITING: 审核中, AUDIT_REJECT: 审核驳回, WAIT_RELEASE: 待上架, BASE_AUDIT_PASS: 准入不可营销, GRAY: 灰度中, RELEASE: 已上架, OFFLINE: 已下架, AUDIT_OFFLINE: 已下架;"`
}

type GetAppVersionDetailReq struct {
	g.Meta     `path:"/:appId/getAppVersionDetail" method:"post" summary:"小程序版本详情查询" tags:"Alipay小程序管理"`
	AppVersion string `json:"app_version" dc:"版本号"`
}
