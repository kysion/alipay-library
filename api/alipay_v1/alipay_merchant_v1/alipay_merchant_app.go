package alipay_merchant_v1

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/kysion/alipay-library/alipay_model"
)

type UpdateStateReq struct {
    g.Meta `path:"/updateState" method:"post" summary:"修改状态" tags:"商家应用"`
    Id     int64 `json:"id" dc:"商家应用Id"`
    State  int   `json:"state" dc:"状态"`
}

type CreateMerchantAppConfigReq struct {
    g.Meta `path:"/createMerchantAppConfig" method:"post" summary:"创建商家应用" tags:"商家应用"`
    alipay_model.AlipayMerchantAppConfig
}

type GetMerchantAppConfigByIdReq struct {
    g.Meta `path:"/getMerchantAppConfigById" method:"post" summary:"根据id获取商家应用" tags:"商家应用"`
    Id     int64 `json:"id" dc:"商家应用Id"`
}

type UpdateMerchantAppConfigReq struct {
    g.Meta `path:"/updateMerchantAppConfig" method:"post" summary:"修改商家应用基础信息" tags:"商家应用"`
    alipay_model.UpdateMerchantAppConfigReq
}

type UpdateMerchantAppConfigHttpsReq struct {
    g.Meta `path:"/updateMerchantAppConfigHttps" method:"post" summary:"修改Https证书认证" tags:"商家应用"`
    alipay_model.UpdateMerchantAppConfigHttpsReq
}

type MerchantAppConfigRes alipay_model.AlipayMerchantAppConfig

type UpdateMerchantKeyCertReq struct {
    g.Meta `path:"/updateMerchantKeyCert" method:"post" summary:"修改密钥证书" tags:"商家应用"`
    alipay_model.UpdateMerchantKeyCertReq
}
