package alipay_third_v1

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/kysion/alipay-test/alipay_model"
)

// 这个文件属于我们调用

// gateway属于微信平台调用或推送

type UpdateStateReq struct {
    g.Meta `path:"/updateState" method:"post" summary:"修改状态" tags:"服务商应用"`
    Id     int64 `json:"id" dc:"服务商应用Id"`
    State  int   `json:"state" dc:"状态"`
}

type CreateThirdAppConfigReq struct {
    g.Meta `path:"/createThirdAppConfig" method:"post" summary:"创建服务商应用并返回开发配置" tags:"服务商应用"`
    alipay_model.AlipayThirdAppConfig
}

type GetThirdAppConfigByIdReq struct {
    g.Meta `path:"/getThirdAppConfigById" method:"post" summary:"根据id获取服务商应用" tags:"服务商应用"`
    Id     int64 `json:"id" dc:"服务商应用Id"`
}

type UpdateThirdAppConfigReq struct {
    g.Meta `path:"/updateThirdAppConfig" method:"post" summary:"修改服务商应用基础信息" tags:"服务商应用"`
    alipay_model.UpdateThirdAppConfig
}

type UpdateThirdAppConfigHttpsReq struct {
    g.Meta `path:"/updateThirdAppConfigHttps" method:"post" summary:"修改Https证书认证" tags:"服务商应用"`
    alipay_model.UpdateThirdAppConfigHttpsReq
}

type UpdateThirdKeyCertReq struct {
    g.Meta `path:"/updateThirdKeyCert" method:"post" summary:"修改密钥证书" tags:"服务商应用"`
    alipay_model.UpdateThirdKeyCertReq
}

type ThirdAppConfigRes alipay_model.AlipayThirdAppConfig
