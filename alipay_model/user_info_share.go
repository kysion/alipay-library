package alipay_model

// UserInfoShareResponse =====================会员信息查询Res==============================
type UserInfoShareResponse struct {
	Response      *UserInfoShare `json:"alipay_user_info_share_response"`
	ErrorResponse *ErrorResponse `json:"error_response,omitempty"`
	AlipayCertSn  string         `json:"alipay_cert_sn,omitempty"`
	SignData      string         `json:"-"`
	Sign          string         `json:"sign"`
}

type UserInfoShare struct {
	UserId             string `json:"user_id,omitempty"  dc:"支付宝唯一标识user_id"`
	Avatar             string `json:"avatar,omitempty"  dc:"头像"`
	Province           string `json:"province,omitempty"  dc:"省份"`
	City               string `json:"city,omitempty"  dc:"城市"`
	NickName           string `json:"nick_name,omitempty"  dc:"昵称"`
	IsStudentCertified string `json:"is_student_certified,omitempty"  dc:"是否学生"`
	UserType           string `json:"user_type,omitempty"  dc:"用户类型"`
	UserStatus         string `json:"user_status,omitempty"  dc:"用户状态"`
	IsCertified        string `json:"is_certified,omitempty"  dc:"是否实名认证"`
	Gender             string `json:"gender,omitempty"  dc:"性别"`
}

type UserInfoRes UserInfoShare
