package alipay_model

type CertifyInitReq struct {
	OuterOrderNo   string         `json:"outer_order_no" dc:"商户请求的唯一标识"`
	BizCode        string         `json:"biz_code" dc:"认证场景码，取值如下:FACE：多因子人脸认证、CERT_PHOTO：多因子证照认证、CERT_PHOTO_FACE ：多因子证照和人脸认证、SMART_FACE：多因子快捷认证"`
	IdentityParam  IdentityParam  `json:"identity_param" dc:"需要验证的身份信息"`
	MerchantConfig MerchantConfig `json:"merchant_config" dc:"商户个性化配置"`
}

type IdentityParam struct {
	IdentityType string `json:"identity_type" dc:"1.若本人验证，使用CERT_INFO；2.若代他人验证，使用AGENT_CERT_INFO；"`
	CertType     string `json:"cert_type" dc:"【可选】身份证: IDENTITY_CARD、港澳居民来往内地通行证: HOME_VISIT_PERMIT_HK_MC、台湾居民来往内地通行证: HOME_VISIT_PERMIT_TAIWAN、港澳居民居住证: RESIDENCE_PERMIT_HK_MC、台湾居民居住证: RESIDENCE_PERMIT_TAIWAN"`
	CertName     string `json:"cert_name" dc:"【可选】填入真实姓名，注意：在identity_type为CERT_INFO或者AGENT_CERT_INFO时，该字段必填"`
	CertNo       string `json:"cert_no" dc:"【可选】填入姓名相匹配的证件号码，注意：在identity_type为CERT_INFO或者AGENT_CERT_INFO时，该字段必填"`
}

type MerchantConfig struct {
	FaceReserveStrategy string `json:"face_reserve_strategy" dc:"【可选】保存活体人脸: reserve、不保存活体人脸: never"`
	ReturnUrl           string `json:"return_url" dc:"认证成功后需要跳转的地址，一般为商户业务页面；若无跳转地址可填空字符"`
}

// ======================================================================================================

type UserCertifyOpenInitRes struct {
	Response     *UserCertifyOpenInit `json:"alipay_user_certify_open_initialize_response"`
	AlipayCertSn string               `json:"alipay_cert_sn,omitempty"`
	SignData     string               `json:"-"`
	Sign         string               `json:"sign"`
}

type UserCertifyOpenInit struct {
	ErrorResponse
	CertifyId string `json:"certify_id,omitempty"`
}

// ======================================================================================================

type UserCertifyOpenQueryRes struct {
	Response     *UserCertifyOpenQuery `json:"alipay_user_certify_open_query_response"`
	AlipayCertSn string                `json:"alipay_cert_sn,omitempty"`
	SignData     string                `json:"-"`
	Sign         string                `json:"sign"`
}

type UserCertifyOpenQuery struct {
	ErrorResponse
	Passed       string `json:"passed,omitempty"`
	IdentityInfo string `json:"identity_info,omitempty"`
	MaterialInfo string `json:"material_info,omitempty"`
}

type UserCertifyOpenRes struct {
	ErrorResponse
	ReturnUrl string `json:"returnUrl"`
	CertifyId string `json:"certify_id,omitempty"`
}
