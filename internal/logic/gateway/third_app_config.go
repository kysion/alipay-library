package gateway

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	dao "github.com/kysion/alipay-library/alipay_model/alipay_dao"
	do "github.com/kysion/alipay-library/alipay_model/alipay_do"
	entity "github.com/kysion/alipay-library/alipay_model/alipay_entity"
	"github.com/kysion/base-library/utility/daoctl"
	"github.com/yitter/idgenerator-go/idgen"
	"strconv"
	"time"
)

type sThirdAppConfig struct {
    redisCache *gcache.Cache
    Duration   time.Duration
}

func NewThirdAppConfig() *sThirdAppConfig {
    return &sThirdAppConfig{
        redisCache: gcache.New(),
    }
}

// GetThirdAppConfigById 根据id查找第三方应用配置信息
func (s *sThirdAppConfig) GetThirdAppConfigById(ctx context.Context, id int64) (*alipay_model.AlipayThirdAppConfig, error) {
    return daoctl.GetByIdWithError[alipay_model.AlipayThirdAppConfig](dao.AlipayThirdAppConfig.Ctx(ctx), id)
}

// GetThirdAppConfigByAppId 根据AppId查找第三方应用配置信息
func (s *sThirdAppConfig) GetThirdAppConfigByAppId(ctx context.Context, id string) (*alipay_model.AlipayThirdAppConfig, error) {
    data := alipay_model.AlipayThirdAppConfig{}

    err := dao.AlipayThirdAppConfig.Ctx(ctx).Where(do.AlipayThirdAppConfig{AppId: id}).Scan(&data)

    return &data, err
}

// GetThirdAppConfigBySysUserId  根据用户id查询第三方应用配置信息
func (s *sThirdAppConfig) GetThirdAppConfigBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayThirdAppConfig, error) {
    result := alipay_model.AlipayThirdAppConfig{}
    model := dao.AlipayThirdAppConfig.Ctx(ctx)

    err := model.Where(dao.AlipayThirdAppConfig.Columns().SysUserId, sysUserId).Scan(&result)

    if err != nil {
        return nil, err
    }

    return &result, nil
}

// CreateThirdAppConfig  创建第三方应用配置信息
func (s *sThirdAppConfig) CreateThirdAppConfig(ctx context.Context, info *alipay_model.AlipayThirdAppConfig) (*alipay_model.AlipayThirdAppConfig, error) {
    // 创建的时候可指定域名，没指定就是用使用当前域名
    // appId的32进制编码
    appId := strconv.FormatInt(gconv.Int64(info.AppId), 36)

    if info.ServerDomain != "" {
        info.AppGatewayUrl = info.ServerDomain + "/merchant/" + appId + "/gateway.services"
        info.AppCallbackUrl = info.ServerDomain + "/merchant/" + appId + "/gateway.callback"
    } else if info.ServerDomain == "" {
        // 没指定服务器域名，默认使用当前服务器域名
        info.ServerDomain = "https://alipay.kuaimk.com"
        info.AppGatewayUrl = "https://alipay.kuaimk.com/alipay/" + appId + "/gateway.services"
        info.AppCallbackUrl = "https://alipay.kuaimk.com/alipay/" + appId + "/gateway.callback"
    }

    // 用户id默认是当前登录用户
    user := sys_service.SysSession().Get(ctx).JwtClaimsUser
    if user.Id != 0 {
        info.SysUserId = user.Id
        info.UnionMainId = user.UnionMainId
    }

    data := do.AlipayThirdAppConfig{}
    gconv.Struct(info, &data)

    data.Id = idgen.NextId()
    data.State = 1 // 状态默认启用
    if info.ExtJson == "" {
        data.ExtJson = nil
    }

    affected, err := daoctl.InsertWithError(
        dao.AlipayThirdAppConfig.Ctx(ctx),
        data,
    )

    if affected == 0 || err != nil {
        return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "第三方应用配置信息创建失败", dao.AlipayThirdAppConfig.Table())
    }

    //  创建好应用了后，平台返回相关开发者配置信息主要就是两个网关地址

    return s.GetThirdAppConfigById(ctx, gconv.Int64(data.Id))
}

// UpdateThirdAppConfig 更新第三方应用基础配置信息
func (s *sThirdAppConfig) UpdateThirdAppConfig(ctx context.Context, info *alipay_model.UpdateThirdAppConfig) (bool, error) {
    // 首先判断第三方应用配置信息是否存在
    consumerInfo, err := daoctl.GetByIdWithError[entity.AlipayThirdAppConfig](dao.AlipayThirdAppConfig.Ctx(ctx), info.Id)
    if err != nil || consumerInfo == nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "该第三方应用配置不存在", dao.AlipayThirdAppConfig.Table())
    }
    data := do.AlipayThirdAppConfig{}
    gconv.Struct(info, &data)
    if info.ExtJson == "" {
        data.ExtJson = nil
    }

    model := dao.AlipayThirdAppConfig.Ctx(ctx)
    affected, err := daoctl.UpdateWithError(model.Data(data).OmitNilData().Where(do.AlipayThirdAppConfig{Id: info.Id}))

    if err != nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "第三方应用配置信息更新失败", dao.AlipayThirdAppConfig.Table())
    }

    return affected > 0, nil
}

// UpdateState 修改第三方应用状态
func (s *sThirdAppConfig) UpdateState(ctx context.Context, id int64, state int) (bool, error) {
    affected, err := daoctl.UpdateWithError(dao.AlipayThirdAppConfig.Ctx(ctx).Data(do.AlipayThirdAppConfig{
        State: state,
    }).OmitNilData().Where(do.AlipayThirdAppConfig{Id: id}))

    if err != nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "第三方应用配置状态修改失败", dao.AlipayThirdAppConfig.Table())
    }
    return affected > 0, err
}

// UpdateAppAuthToken 更新Token  服务商应用授权token
func (s *sThirdAppConfig) UpdateAppAuthToken(ctx context.Context, info *alipay_model.UpdateThirdAppAuthToken) (bool, error) {
    data := do.AlipayThirdAppConfig{}
    gconv.Struct(info, &data)

    affected, err := daoctl.UpdateWithError(dao.AlipayThirdAppConfig.Ctx(ctx).Data(data).OmitNilData().Where(do.AlipayThirdAppConfig{AppId: info.AppId}))

    if err != nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "服务商应用Token修改失败", dao.AlipayThirdAppConfig.Table())
    }
    return affected > 0, err
}

// UpdateAppConfigHttps 修改服务商应用Https配置
func (s *sThirdAppConfig) UpdateAppConfigHttps(ctx context.Context, info *alipay_model.UpdateThirdAppConfigHttpsReq) (bool, error) {
    data := do.AlipayThirdAppConfig{}
    gconv.Struct(info, &data)

    affected, err := daoctl.UpdateWithError(dao.AlipayThirdAppConfig.Ctx(ctx).Data(data).OmitNilData().Where(do.AlipayThirdAppConfig{Id: info.Id}))

    if err != nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "服务商应用基础修改失败", dao.AlipayThirdAppConfig.Table())
    }
    return affected > 0, err
}

// UpdateThirdKeyCert 更新第三方应用配置证书密钥
func (s *sThirdAppConfig) UpdateThirdKeyCert(ctx context.Context, info *alipay_model.UpdateThirdKeyCertReq) (bool, error) {
    app, err := s.GetThirdAppConfigByAppId(ctx, info.AppId)
    if err != nil || app == nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "该第三方应用配置不存在", dao.AlipayThirdAppConfig.Table())
    }

    data := do.AlipayThirdAppConfig{}
    gconv.Struct(info, &data)

    model := dao.AlipayThirdAppConfig.Ctx(ctx)
    affected, err := daoctl.UpdateWithError(model.Data(data).OmitNilData().Where(do.AlipayThirdAppConfig{AppId: info.AppId}))

    if err != nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "第三方应用密钥证书更新失败", dao.AlipayThirdAppConfig.Table())
    }

    return affected > 0, nil
}
