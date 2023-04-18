package app_version_manager

import "github.com/kysion/base-library/utility/enum"

// VersionStatusEnum 版本状态：INIT: 开发中, AUDITING: 审核中, AUDIT_REJECT: 审核驳回, WAIT_RELEASE: 待上架, BASE_AUDIT_PASS: 准入不可营销, GRAY: 灰度中, RELEASE: 已上架, OFFLINE: 已下架, AUDIT_OFFLINE: 已下架;
type VersionStatusEnum enum.IEnumCode[string]

type versionStatus struct {
	INIT            VersionStatusEnum
	AUDITING        VersionStatusEnum
	AUDIT_REJECT    VersionStatusEnum
	WAIT_RELEASE    VersionStatusEnum
	BASE_AUDIT_PASS VersionStatusEnum
	GRAY            VersionStatusEnum
	RELEASE         VersionStatusEnum
	OFFLINE         VersionStatusEnum
}

var Status = versionStatus{
	INIT:            enum.New[VersionStatusEnum]("INIT", "开发中"),
	AUDITING:        enum.New[VersionStatusEnum]("AUDITING", "审核中"),
	AUDIT_REJECT:    enum.New[VersionStatusEnum]("AUDIT_REJECT", "审核驳回"),
	WAIT_RELEASE:    enum.New[VersionStatusEnum]("WAIT_RELEASE", "待上架"),
	BASE_AUDIT_PASS: enum.New[VersionStatusEnum]("BASE_AUDIT_PASS", "准入不可营销"),
	GRAY:            enum.New[VersionStatusEnum]("GRAY", "灰度中"),
	RELEASE:         enum.New[VersionStatusEnum]("RELEASE", "已上架"),
	OFFLINE:         enum.New[VersionStatusEnum]("OFFLINE", "已下架"),
}
