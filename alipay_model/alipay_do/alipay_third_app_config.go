// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package share_do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AlipayThirdAppConfig is the golang structure of table alipay_third_app_config for DAO operations like Where/Data.
type AlipayThirdAppConfig struct {
	g.Meta                  `orm:"table:alipay_third_app_config, do:true"`
	Id                      interface{} // 服务商id
	Name                    interface{} // 服务商name
	AppId                   interface{} // 服务商应用Id
	AppName                 interface{} // 服务商应用名称
	AppType                 interface{} // 服务商应用类型：1小程序  2网站/移动应用  4生活号
	AppAuthToken            interface{} // 服务商授权应用token
	State                   interface{} // 状态：0禁用 1启用
	ExpiresIn               *gtime.Time // Token失效时间
	ReExpiresIn             *gtime.Time // Token刷新限期时间
	UserId                  interface{} // 应用所属账号
	UnionMainId             interface{} // 关联主体id
	SysUserId               interface{} // 用户id
	ExtJson                 interface{} // 拓展字段
	AppGatewayUrl           interface{} // 网关地址
	AppCallbackUrl          interface{} // 回调地址
	AesEncryptKey           interface{} // AES接口内容加密方式
	ServerDomain            interface{} // 服务器域名
	Logo                    interface{} // 服务商应用logo
	HttpsCert               interface{} // 域名证书
	HttpsKey                interface{} // 域名私钥
	PrivateKey              interface{} // 私钥
	PublicKey               interface{} // 公钥
	PublicKeyCert           interface{} // 公钥证书
	AppPublicCertKey        interface{} // 应用证书公钥
	AlipayRootCertPublicKey interface{} // 阿里根证书公钥
	DevState                interface{} // 开发状态：0未上线 1已上线
	InterfaceSignType       interface{} // 接口加签方式：1密钥 2证书
	UpdatedAt               *gtime.Time //
	AppIdMd5                interface{} // 应用id加密md5后的结果
}
