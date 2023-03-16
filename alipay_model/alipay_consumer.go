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
