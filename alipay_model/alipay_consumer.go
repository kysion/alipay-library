package alipay_model

import "github.com/gogf/gf/v2/os/gtime"

type AlipayConsumerConfig struct {
	Id                 int64       `json:"id"                 dc:"id"`
	UserId             string      `json:"userId"             dc:"用户账号id"`
	SysUserId          int64       `json:"sysUserId"          dc:"用户id"`
	Avatar             string      `json:"avatar"             dc:"头像"`
	Province           string      `json:"province"           dc:"省份"`
	City               string      `json:"city"               dc:"城市"`
	NickName           string      `json:"nickName"           dc:"昵称"`
	IsStudentCertified int         `json:"isStudentCertified" dc:"学生认证"`
	UserType           string      `json:"userType"           dc:"用户账号类型"`
	UserState          int         `json:"userState"          dc:"用户状态"`
	IsCertified        int         `json:"isCertified"        dc:"是否实名认证"`
	Sex                int         `json:"sex"                dc:"性别：0未知、1男、2女"`
	AuthToken          string      `json:"authToken"          dc:"授权token"`
	ExtJson            string      `json:"extJson"            dc:"拓展字段"`
	AuthState          int         `json:"authState"          description:"用户授权状态：1授权、2未授权"`
	AlipayUserId       string      `json:"alipayUserId"       description:"Alipay的UserId"`
	ExpiresIn          *gtime.Time `json:"expiresIn"          description:"用户授权Token过期时间"`
	ReFreshToken       string      `json:"reFreshToken"       description:"刷新Token"`
	ReExpiresIn        *gtime.Time `json:"reExpiresIn"        description:"刷新Token过期时间"`
	AuthStart          *gtime.Time `json:"authStart"          description:"用户授权开始时间"`
	AppType            int         `json:"appType"            description:"应用类型：1小程序  2网站/移动应用  4生活号"`
	AppId              string      `json:"appId"              description:"商家应用Id"`
}

type UpdateConsumerReq struct {
	Id                 *int64      `json:"id"                 dc:"id"`
	Avatar             *string     `json:"avatar"             dc:"头像"`
	Province           *string     `json:"province"           dc:"省份"`
	City               *string     `json:"city"               dc:"城市"`
	NickName           *string     `json:"nickName"           dc:"昵称"`
	IsStudentCertified *int        `json:"isStudentCertified" dc:"学生认证"`
	AuthToken          *string     `json:"authToken"          dc:"授权token"`
	ExtJson            *string     `json:"extJson"            dc:"拓展字段"`
	AuthState          *int        `json:"authState"          description:"用户授权状态：1授权、2未授权"`
	AlipayUserId       *string     `json:"alipayUserId"       description:"Alipay的UserId"`
	ExpiresIn          *gtime.Time `json:"expiresIn"          description:"用户授权Token过期时间"`
	ReFreshToken       *string     `json:"reFreshToken"       description:"刷新Token"`
	ReExpiresIn        *gtime.Time `json:"reExpiresIn"        description:"刷新Token过期时间"`
	AuthStart          *gtime.Time `json:"authStart"          description:"用户授权开始时间"`
	AppType            *int        `json:"appType"            description:"应用类型：1小程序  2网站/移动应用  4生活号"`
	AppId              *string     `json:"appId"              description:"商家应用Id"`
}
