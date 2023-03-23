package alipay_model

import "github.com/gogf/gf/v2/os/gtime"

type AlipayMerchantAppConfig struct {
	Id                  int64       `json:"id"                  description:"商家id"`
	Name                string      `json:"name"                description:"商家name"`
	AppId               string      `json:"appId"               description:"商家应用Id"`
	AppName             string      `json:"appName"             description:"商家应用名称"`
	AppType             int         `json:"appType"             description:"应用类型：1小程序  2网站/移动应用  4生活号"`
	AppAuthToken        string      `json:"appAuthToken"        description:"商家授权应用token"`
	IsFullProxy         int         `json:"isFullProxy"         description:"是否全权委托待开发：0否 1是"`
	State               int         `json:"state"               description:"状态：0禁用 1启用"`
	ExpiresIn           *gtime.Time `json:"expiresIn"           description:"Token失效时间"`
	ReExpiresIn         *gtime.Time `json:"reExpiresIn"         description:"Token刷新限期时间"`
	UserId              string      `json:"userId"              description:"应用所属账号"`
	UnionMainId         int64       `json:"unionMainId"         description:"关联主体id"`
	SysUserId           int64       `json:"sysUserId"           description:"用户id"`
	ExtJson             string      `json:"extJson"             description:"拓展字段"`
	AppGatewayUrl       string      `json:"appGatewayUrl"       description:"网关地址"`
	AppCallbackUrl      string      `json:"appCallbackUrl"      description:"回调地址"`
	AesEncryptKey       string      `json:"aesEncryptKey"       description:"AES接口内容加密方式"`
	ServerDomain        string      `json:"serverDomain"        description:"服务器域名"`
	Logo                string      `json:"logo"                description:"商家应用logo"`
	HttpsCert           string      `json:"httpsCert"           description:"域名证书"`
	HttpsKey            string      `json:"httpsKey"            description:"域名私钥"`
	PrivateKey          string      `json:"privateKey"          description:"私钥"`
	PublicKey           string      `json:"publicKey"           description:"公钥"`
	PublicKeyCert       string      `json:"publicKeyCert"       description:"公钥证书"`
	AppPublicCertKey    string      `json:"appPublicCertKey"    description:"应用证书公钥"`
	AlipayCertPublicKey string      `json:"alipayCertPublicKey" description:"阿里根证书公钥"`
	DevState            int         `json:"devState"            description:"开发状态：0未上线 1已上线"`
	InterfaceSignType   int         `json:"interfaceSignType"   description:"接口加签方式：1密钥 2证书"`
	AppIdMd5            string      `json:"appIdMd5"            description:"应用id加密md5后的结果"`
	ThirdAppId          string      `json:"thirdAppId"              description:"服务商appId"`
	NotifyUrl           string      `json:"notifyUrl"               description:"异步通知地址，允许业务层追加相关参数"`
}

// UpdateMerchantAppConfigReq 修改商户基础信息
type UpdateMerchantAppConfigReq struct {
	Id             int64       `json:"id"                  description:"商家id"`
	Name           string      `json:"name"                description:"商家name"`
	AppAuthToken   string      `json:"appAuthToken"        description:"商家授权应用token"`
	ExpiresIn      *gtime.Time `json:"expiresIn"           description:"Token失效时间"`
	ReExpiresIn    *gtime.Time `json:"reExpiresIn"         description:"Token刷新限期时间"`
	ExtJson        string      `json:"extJson"             description:"拓展字段"`
	AppGatewayUrl  string      `json:"appGatewayUrl"       description:"网关地址"`
	AppCallbackUrl string      `json:"appCallbackUrl"      description:"回调地址"`
	AesEncryptKey  string      `json:"aesEncryptKey"       description:"AES接口内容加密方式"`
	ServerDomain   string      `json:"serverDomain"        description:"服务器域名"`
	Logo           string      `json:"logo"                description:"商家应用logo"`
	AppIdMd5       string      `json:"appIdMd5"            description:"应用id加密md5后的结果"`
	NotifyUrl      string      `json:"notifyUrl"               description:"异步通知地址，允许业务层追加相关参数"`
}

type UpdateMerchantAppAuthToken struct {
	AppId        string      `json:"appId"               description:"商家应用Id"`
	AppAuthToken string      `json:"appAuthToken"        description:"商家授权应用token"`
	ExpiresIn    *gtime.Time `json:"expiresIn"           description:"Token失效时间"`
	ReExpiresIn  *gtime.Time `json:"reExpiresIn"         description:"Token刷新限期时间"`
	ThirdAppId   string      `json:"thirdAppId"              description:"服务商appId"`
}

// UpdateMerchantAppConfigHttpsReq 修改https证书
type UpdateMerchantAppConfigHttpsReq struct {
	Id        int64  `json:"id"                  description:"商家id"`
	HttpsCert string `json:"httpsCert"           description:"域名证书"`
	HttpsKey  string `json:"httpsKey"            description:"域名私钥"`
}

// UpdateMerchantKeyCertReq 修改证书公钥或者密钥
type UpdateMerchantKeyCertReq struct {
	AppId               string `json:"appId"               description:"商家应用Id"`
	PrivateKey          string `json:"privateKey"          description:"私钥"`
	PublicKey           string `json:"publicKey"           description:"公钥"`
	PublicKeyCert       string `json:"publicKeyCert"       description:"公钥证书"`
	AppPublicCertKey    string `json:"appPublicCertKey"    description:"应用证书公钥"`
	AlipayCertPublicKey string `json:"alipayCertPublicKey" description:"阿里根证书公钥"`
}
