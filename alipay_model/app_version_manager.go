package alipay_model

// V2版本公共响应参数------------------------------------------------------------------------------------------------------------------------------------------------

type SubmitAppVersionAuditRes struct {
	Sign                                    string                                  `json:"sign"`
	AlipayOpenMiniVersionAuditApplyResponse AlipayOpenMiniVersionAuditApplyResponse `json:"alipay_open_mini_version_audit_apply_response"`
}

type AlipayOpenMiniVersionAuditApplyResponse struct {
	ErrorResponse
	SpeedUp     string `json:"speed_up"`
	SpeedUpMemo string `json:"speed_up_memo"`
}

// V3版本公共响应参数------------------------------------------------------------------------------------------------------------------------------------------------

type Common struct {
	AlipayTimestamp string `json:"alipay_timestamp" dc:"时间"`
	AlipaySignature string `json:"alipay_signature" dc:"签名"`
	AlipayTraceid   string `json:"alipay_traceid" dc:"支付宝traceId ，用于排查问题使用"`
	AlipayNonce     string `json:"alipay_nonce" dc:"支付宝nonce标记，每次请求会生成不同的nonce，可用于防重放判断"`
}

// 小程序提交审核------------------------------------------------------------------------------------------------------------------------------------------------

type AppVersionAuditReq struct {
	AppVersion   string `json:"app_version" dc:"应用版本号"`
	VersionDesc  string `json:"version_desc" dc:"版本描述"`
	RegionType   string `json:"region_type" dc:"小程序服务区域类型:支持GLOBAL-全球 CHINA-中国 如果前期已经设置过该信息，本次可不填"`
	SpeedUp      string `json:"speed_up" dc:"是否加速审核，需要拥有绿通权益"`
	AutoOnline   string `json:"auto_online" dc:"是否自动上线"`
	TestFileName string `json:"test_file_name" dc:"测试附件，用于上传测试报告和测试录屏，请上传10M以内附件，支持格式zip，rar"`
	TestAccout   string `json:"test_accout" dc:"测试账号"`
	TestPassword string `json:"test_password" dc:"测试账号密码"`
}

type AppVersionAuditPicReq struct {
	FirstLicensePicPath         string `json:"first_license_pic_path,omitempty"`
	SecondLicensePicPath        string `json:"second_license_pic_path,omitempty"`
	ThirdLicensePicPath         string `json:"third_license_pic_path,omitempty"`
	FourthLicensePicPath        string `json:"fourth_license_pic_path,omitempty"`
	FifthLicensePicPath         string `json:"fifth_license_pic_path,omitempty"`
	FirstScreenShotPath         string `json:"first_screen_shot_path,omitempty"`
	SecondScreenShotPath        string `json:"second_screen_shot_path,omitempty"`
	ThirdScreenShotPath         string `json:"third_screen_shot_path,omitempty"`
	FourthScreenShotPath        string `json:"fourth_screen_shot_path,omitempty"`
	FifthScreenShotPath         string `json:"fifth_screen_shot_path,omitempty"`
	FirstSpecialLicensePicPath  string `json:"first_special_license_pic_path,omitempty"`
	SecondSpecialLicensePicPath string `json:"second_special_license_pic_path,omitempty"`
	ThirdSpecialLicensePicPath  string `json:"third_special_license_pic_path,omitempty"`
}

type AppVersionAuditRes struct {
	Common
	SpeedUp     string `json:"speed_up" dc:"是否加速审核"`
	SpeedUpMemo string `json:"speed_up_memo" dc:"提审加速审核说明"`
}

// 小程序撤销审核--------------------------------------------------------------------------------------------------

type CancelVersionAuditRes struct {
	Common
}

// 小程序退回开发版本--------------------------------------------------------------------------------------------------

type CancelVersionRes struct {
	Common
}

// 小程序上架--------------------------------------------------------------------------------------------------

type AppOnlineRes struct {
	Common
}

// 小程序下架--------------------------------------------------------------------------------------------------

type AppOfflineRes struct {
	Common
}

// 获取小程序版本列表--------------------------------------------------------------------------------------------------

type QueryAppVersionListRes struct {
	AppVersionInfos []AppVersionInfos `json:"app_version_infos"`
	AppVersions     []string          `json:"app_versions"`
}

type AppVersionInfos struct {
	BundleID           string `json:"bundle_id"`
	AppVersion         string `json:"app_version"`
	VersionDescription string `json:"version_description"`
	VersionStatus      string `json:"version_status"`
	CreateTime         string `json:"create_time"`
	BaseAudit          string `json:"base_audit"`
	PromoteAudit       string `json:"promote_audit"`
	CanRelease         string `json:"can_release"`
}

// 获取小程序版本详情--------------------------------------------------------------------------------------------------

type QueryAppVersionDetailRes struct {
	AppVersion              string                    `json:"app_version"`
	AppName                 string                    `json:"app_name"`
	AppEnglishName          string                    `json:"app_english_name"`
	AppLogo                 string                    `json:"app_logo"`
	VersionDesc             string                    `json:"version_desc"`
	GrayStrategy            string                    `json:"gray_strategy"`
	Status                  string                    `json:"status"`
	RejectReason            string                    `json:"reject_reason"`
	ScanResult              string                    `json:"scan_result"`
	GmtCreate               string                    `json:"gmt_create"`
	GmtApplyAudit           string                    `json:"gmt_apply_audit"`
	GmtOnline               string                    `json:"gmt_online"`
	GmtOffline              string                    `json:"gmt_offline"`
	AppDesc                 string                    `json:"app_desc"`
	GmtAuditEnd             string                    `json:"gmt_audit_end"`
	ServiceRegionType       string                    `json:"service_region_type"`
	ServiceRegionInfo       []ServiceRegionInfo       `json:"service_region_info"`
	ScreenShotList          []string                  `json:"screen_shot_list"`
	AppSlogan               string                    `json:"app_slogan"`
	Memo                    string                    `json:"memo"`
	ServicePhone            string                    `json:"service_phone"`
	ServiceEmail            string                    `json:"service_email"`
	MiniAppCategoryInfoList []MiniAppCategoryInfoList `json:"mini_app_category_info_list"`
	PackageInfoList         []PackageInfoList         `json:"package_info_list"`
	MiniCategoryInfoList    []MiniCategoryInfoList    `json:"mini_category_info_list"`
	BaseAudit               string                    `json:"base_audit"`
	PromoteAudit            string                    `json:"promote_audit"`
	CanRelease              string                    `json:"can_release"`
	BaseAuditRecord         BaseAuditRecord           `json:"base_audit_record"`
	PromoteAuditRecord      PromoteAuditRecord        `json:"promote_audit_record"`
}
type ServiceRegionInfo struct { // ServiceRegionInfo 省市区信息，当区域类型为LOCATION时，不能为空
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
	CityCode     string `json:"city_code"`
	CityName     string `json:"city_name"`
	AreaCode     string `json:"area_code"`
	AreaName     string `json:"area_name"`
}
type MiniAppCategoryInfoList struct {
	FirstCategoryID    string `json:"first_category_id"`
	FirstCategoryName  string `json:"first_category_name"`
	SecondCategoryID   string `json:"second_category_id"`
	SecondCategoryName string `json:"second_category_name"`
	ThirdCategoryID    string `json:"third_category_id"`
	ThirdCategoryName  string `json:"third_category_name"`
}
type PackageInfoList struct {
	PackageName     string `json:"package_name"`
	PackageDesc     string `json:"package_desc"`
	DocURL          string `json:"doc_url"`
	Status          string `json:"status"`
	PackageOpenType string `json:"package_open_type"`
}
type MiniCategoryInfoList struct {
	FirstCategoryID    string `json:"first_category_id"`
	FirstCategoryName  string `json:"first_category_name"`
	SecondCategoryID   string `json:"second_category_id"`
	SecondCategoryName string `json:"second_category_name"`
	ThirdCategoryID    string `json:"third_category_id"`
	ThirdCategoryName  string `json:"third_category_name"`
}
type Memos struct {
	Memo          string   `json:"memo"`
	MemoImageList []string `json:"memo_image_list"`
}
type BaseAuditRecord struct {
	AuditImages []string `json:"audit_images"`
	Memos       []Memos  `json:"memos"`
}
type PromoteAuditRecord struct {
	AuditImages []string `json:"audit_images"`
	Memos       []Memos  `json:"memos"`
}
