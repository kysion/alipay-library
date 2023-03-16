package alipay_model

type UserInfoShareResponse struct {
	Response      *UserInfoShare `json:"alipay_user_info_share_response"`
	ErrorResponse *ErrorResponse `json:"error_response,omitempty"`
	AlipayCertSn  string         `json:"alipay_cert_sn,omitempty"`
	SignData      string         `json:"-"`
	Sign          string         `json:"sign"`
}

type UserInfoShare struct {
	UserId             string `json:"user_id,omitempty"`
	Avatar             string `json:"avatar,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	NickName           string `json:"nick_name,omitempty"`
	IsStudentCertified string `json:"is_student_certified,omitempty"`
	UserType           string `json:"user_type,omitempty"`
	UserStatus         string `json:"user_status,omitempty"`
	IsCertified        string `json:"is_certified,omitempty"`
	Gender             string `json:"gender,omitempty"`
}

type UserInfoRes UserInfoShare
