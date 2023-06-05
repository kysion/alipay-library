// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package alipay_entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AlipayThirdAppConfig is the golang structure for table alipay_third_app_config.
type AlipayThirdAppConfig struct {
	Id                      int64       `json:"id"                      description:"服务商id"`
	Name                    string      `json:"name"                    description:"服务商name"`
	AppId                   string      `json:"appId"                   description:"服务商应用Id"`
	AppName                 string      `json:"appName"                 description:"服务商应用名称"`
	AppType                 int         `json:"appType"                 description:"服务商应用类型：1小程序  2网站/移动应用  4生活号"`
	AppAuthToken            string      `json:"appAuthToken"            description:"服务商授权应用token"`
	State                   int         `json:"state"                   description:"状态：0禁用 1启用"`
	ExpiresIn               *gtime.Time `json:"expiresIn"               description:"Token失效时间"`
	ReExpiresIn             *gtime.Time `json:"reExpiresIn"             description:"Token刷新限期时间"`
	UserId                  string      `json:"userId"                  description:"应用所属账号"`
	UnionMainId             int64       `json:"unionMainId"             description:"关联主体id"`
	SysUserId               int64       `json:"sysUserId"               description:"用户id"`
	ExtJson                 string      `json:"extJson"                 description:"拓展字段"`
	AppGatewayUrl           string      `json:"appGatewayUrl"           description:"网关地址"`
	AppCallbackUrl          string      `json:"appCallbackUrl"          description:"回调地址"`
	AesEncryptKey           string      `json:"aesEncryptKey"           description:"AES接口内容加密方式"`
	ServerDomain            string      `json:"serverDomain"            description:"服务器域名"`
	Logo                    string      `json:"logo"                    description:"服务商应用logo"`
	HttpsCert               string      `json:"httpsCert"               description:"域名证书"`
	HttpsKey                string      `json:"httpsKey"                description:"域名私钥"`
	PrivateKey              string      `json:"privateKey"              description:"私钥"`
	PublicKey               string      `json:"publicKey"               description:"公钥"`
	PublicKeyCert           string      `json:"publicKeyCert"           description:"公钥证书"`
	AppPublicCertKey        string      `json:"appPublicCertKey"        description:"应用证书公钥"`
	AlipayRootCertPublicKey string      `json:"alipayRootCertPublicKey" description:"阿里根证书公钥"`
	DevState                int         `json:"devState"                description:"开发状态：0未上线 1已上线"`
	InterfaceSignType       int         `json:"interfaceSignType"       description:"接口加签方式：1密钥 2证书"`
	UpdatedAt               *gtime.Time `json:"updatedAt"               description:""`
	AppIdMd5                string      `json:"appIdMd5"                description:"应用id加密md5后的结果"`
}
