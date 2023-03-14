package alipay_consts

type global struct {
	RSA2                    string
	PriPath                 string
	PublicCrtPath           string
	AppCertPublicKeyPath    string
	AlipayRootCertPath      string
	AlipayCertPublicKeyPath string
	AppId                   string
	AppCode                 string
	AES                     string
	CallbackUrl             string
	ReturnUrl               string
	TradeHookExpireAt       int64
}

var (
	Global = global{}
)
