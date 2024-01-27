package alipay_model

type AlipayConsumerConfig struct {
	Id                 int64  `json:"id"                 dc:"id"`
	UserId             string `json:"userId"             dc:"用户账号id"`
	SysUserId          int64  `json:"sysUserId"          dc:"用户id"`
	Avatar             string `json:"avatar"             dc:"头像"`
	Province           string `json:"province"           dc:"省份"`
	City               string `json:"city"               dc:"城市"`
	NickName           string `json:"nickName"           dc:"昵称"`
	IsStudentCertified int    `json:"isStudentCertified" dc:"学生认证"`
	UserType           string `json:"userType"           dc:"用户账号类型"`
	UserState          int    `json:"userState"          dc:"用户状态"`
	IsCertified        int    `json:"isCertified"        dc:"是否实名认证"`
	Sex                int    `json:"sex"                dc:"性别：0女 1男"`
	AuthToken          string `json:"authToken"          dc:"授权token"`
	ExtJson            string `json:"extJson"            dc:"拓展字段"`
}

type UpdateConsumerReq struct {
	Id                 int64  `json:"id"                 dc:"id"`
	Avatar             string `json:"avatar"             dc:"头像"`
	Province           string `json:"province"           dc:"省份"`
	City               string `json:"city"               dc:"城市"`
	NickName           string `json:"nickName"           dc:"昵称"`
	IsStudentCertified int    `json:"isStudentCertified" dc:"学生认证"`
	AuthToken          string `json:"authToken"          dc:"授权token"`
	ExtJson            string `json:"extJson"            dc:"拓展字段"`
}

type OpenInitialize struct {
	OuterOrderNo  string `json:"outerOrderNo " dc:"商户请求的唯一的标识符"`
	BizCode       string `json:"bizCode" dc:"认证场景码"`
	IdentityParam IdentityParam
}

type ResponseRes struct {
	ErrorResponse
	CertifyId string `json:"certify_id" dc:"本次申请操作的唯一标识"`
}

type IdentityParam struct {
	IdentityType string `json:"identityType" dc:"身份信息参数类型" default:"CERT_INFO"`
	CertType     string `json:"certType" dc:"证件类型" v:"required#请选择证件类型"`
	CertName     string `json:"certName" dc:"真实姓名" v:"required#请输入真实名字"`
	CertNo       string `json:"certNo" dc:"证件号码" v:"required#请输入证件号码"`
}

type MerchantConfig struct {
	ReturnUrl string `json:"returnUrl" dc:"需要回跳的目标地址"`
}
