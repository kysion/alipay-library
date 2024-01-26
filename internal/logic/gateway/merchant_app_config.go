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
	"github.com/kysion/alipay-library/alipay_utility"
	"github.com/kysion/base-library/utility/daoctl"
	"github.com/yitter/idgenerator-go/idgen"
	"strconv"
	"time"
)

type sMerchantAppConfig struct {
	redisCache *gcache.Cache
	Duration   time.Duration
}

func NewMerchantAppConfig() *sMerchantAppConfig {
	return &sMerchantAppConfig{
		redisCache: gcache.New(),
	}
}

// GetMerchantAppConfigById 根据id查找商家应用配置信息
func (s *sMerchantAppConfig) GetMerchantAppConfigById(ctx context.Context, id int64) (*alipay_model.AlipayMerchantAppConfig, error) {
	return daoctl.GetByIdWithError[alipay_model.AlipayMerchantAppConfig](dao.AlipayMerchantAppConfig.Ctx(ctx), id)
}

// GetMerchantAppConfigByAppId 根据AppId查找商家应用配置信息
func (s *sMerchantAppConfig) GetMerchantAppConfigByAppId(ctx context.Context, id string) (*alipay_model.AlipayMerchantAppConfig, error) {
	var data *alipay_model.AlipayMerchantAppConfig

	err := dao.AlipayMerchantAppConfig.Ctx(ctx).Where(do.AlipayMerchantAppConfig{AppId: id}).Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, err
}

// GetMerchantAppConfigBySysUserId  根据商家id查询商家应用配置信息
func (s *sMerchantAppConfig) GetMerchantAppConfigBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayMerchantAppConfig, error) {
	result := alipay_model.AlipayMerchantAppConfig{}
	model := dao.AlipayMerchantAppConfig.Ctx(ctx)

	err := model.Where(dao.AlipayMerchantAppConfig.Columns().SysUserId, sysUserId).Scan(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateMerchantAppConfig  创建商家应用配置信息
func (s *sMerchantAppConfig) CreateMerchantAppConfig(ctx context.Context, info *alipay_model.AlipayMerchantAppConfig) (*alipay_model.AlipayMerchantAppConfig, error) {
	// 创建的时候可指定域名，没指定就是用使用当前域名
	// appId的32进制编码
	appId := strconv.FormatInt(gconv.Int64(info.AppId), 32)

	if info.ServerDomain != "" {
		//appIdHash := alipay_utility.Md5Hash(info.AppId)
		//// 取其appId Md5加密后的前16位  //https://alipay.jditco.com/alipay/appIdMd5-16/gateway.services
		info.AppGatewayUrl = info.ServerDomain + "/alipay/" + appId + "/gateway.services"
		info.AppCallbackUrl = info.ServerDomain + "/alipay/" + appId + "/gateway.callback"
		info.NotifyUrl = info.ServerDomain + "/alipay/" + appId + "/gateway.notify"

		//info.AppIdMd5 = appIdHash
	} else if info.ServerDomain == "" {
		// 没指定服务器域名，默认使用当前服务器域名
		info.ServerDomain = "https://www.kuaimk.com"
		info.AppGatewayUrl = "https://www.kuaimk.com/alipay/" + appId + "/gateway.services"
		info.AppCallbackUrl = "https://www.kuaimk.com/alipay/" + appId + "/gateway.callback"
		info.NotifyUrl = "https://www.kuaimk.com/alipay/" + appId + "/gateway.notify"
	}

	// 用户id默认是当前登录用户
	user := sys_service.SysSession().Get(ctx).JwtClaimsUser
	if user.Id != 0 {
		info.SysUserId = user.Id
		info.UnionMainId = user.UnionMainId
	}

	data := do.AlipayMerchantAppConfig{}

	gconv.Struct(info, &data)

	data.Id = idgen.NextId()
	data.State = 1 // 状态默认启用
	if info.ExtJson == "" {
		data.ExtJson = nil
	}

	affected, err := daoctl.InsertWithError(
		dao.AlipayMerchantAppConfig.Ctx(ctx),
		data,
	)

	if affected == 0 || err != nil {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "商家应用配置信息创建失败", dao.AlipayMerchantAppConfig.Table())
	}

	return s.GetMerchantAppConfigById(ctx, gconv.Int64(data.Id))
}

// UpdateMerchantAppConfig 更新商家应用配置信息
func (s *sMerchantAppConfig) UpdateMerchantAppConfig(ctx context.Context, info *alipay_model.UpdateMerchantAppConfigReq) (bool, error) {
	// 首先判断商家应用配置信息是否存在
	consumerInfo, err := daoctl.GetByIdWithError[entity.AlipayMerchantAppConfig](dao.AlipayMerchantAppConfig.Ctx(ctx), info.Id)
	if err != nil || consumerInfo == nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "该商家应用配置不存在", dao.AlipayMerchantAppConfig.Table())
	}
	data := do.AlipayMerchantAppConfig{}
	gconv.Struct(info, &data)
	if info.ExtJson == "" {
		data.ExtJson = nil
	}
	model := dao.AlipayMerchantAppConfig.Ctx(ctx)
	affected, err := daoctl.UpdateWithError(model.Data(data).OmitNilData().Where(do.AlipayMerchantAppConfig{Id: info.Id}))

	if err != nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "商家应用配置信息更新失败", dao.AlipayMerchantAppConfig.Table())
	}

	return affected > 0, nil
}

// UpdateState 修改商家应用状态
func (s *sMerchantAppConfig) UpdateState(ctx context.Context, id int64, state int) (bool, error) {
	affected, err := daoctl.UpdateWithError(dao.AlipayMerchantAppConfig.Ctx(ctx).Data(do.AlipayMerchantAppConfig{
		State: state,
	}).OmitNilData().Where(do.AlipayMerchantAppConfig{Id: id}))

	if err != nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "商家应用状态修改失败", dao.AlipayMerchantAppConfig.Table())
	}
	return affected > 0, err
}

// UpdateAppAuthToken 更新Token  商家应用授权token
func (s *sMerchantAppConfig) UpdateAppAuthToken(ctx context.Context, info *alipay_model.UpdateMerchantAppAuthToken) (bool, error) {
	data := do.AlipayMerchantAppConfig{}
	gconv.Struct(info, &data)

	affected, err := daoctl.UpdateWithError(dao.AlipayMerchantAppConfig.Ctx(ctx).Data(data).OmitNilData().Where(do.AlipayMerchantAppConfig{AppId: info.AppId}))

	if err != nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "商家应用Token修改失败", dao.AlipayMerchantAppConfig.Table())
	}
	return affected > 0, err
}

// UpdateAppConfigHttps 修改商家应用Https配置
func (s *sMerchantAppConfig) UpdateAppConfigHttps(ctx context.Context, info *alipay_model.UpdateMerchantAppConfigHttpsReq) (bool, error) {
	data := do.AlipayMerchantAppConfig{}
	gconv.Struct(info, &data)

	affected, err := daoctl.UpdateWithError(dao.AlipayMerchantAppConfig.Ctx(ctx).Data(data).OmitNilData().Where(do.AlipayMerchantAppConfig{Id: info.Id}))

	if err != nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "商家应用基础修改失败", dao.AlipayMerchantAppConfig.Table())
	}
	return affected > 0, err
}

// UpdateMerchantKeyCert 更新商家应用配置证书密钥
func (s *sMerchantAppConfig) UpdateMerchantKeyCert(ctx context.Context, info *alipay_model.UpdateMerchantKeyCertReq) (bool, error) {
	app, err := s.GetMerchantAppConfigByAppId(ctx, info.AppId)
	if err != nil || app == nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "该商家应用配置不存在", dao.AlipayMerchantAppConfig.Table())
	}

	data := do.AlipayMerchantAppConfig{}
	gconv.Struct(info, &data)

	model := dao.AlipayMerchantAppConfig.Ctx(ctx)
	affected, err := daoctl.UpdateWithError(model.Data(data).OmitNilData().Where(do.AlipayMerchantAppConfig{AppId: info.AppId}))

	if err != nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "商家应用密钥证书更新失败", dao.AlipayMerchantAppConfig.Table())
	}

	return affected > 0, nil
}

// CreatePolicy 创建用户协议或隐私协议
func (s *sMerchantAppConfig) CreatePolicy(ctx context.Context, info *alipay_model.CreatePolicyReq) (bool, error) {
	appId := alipay_utility.GetAlipayAppIdFormCtx(ctx)

	app, err := s.GetMerchantAppConfigByAppId(ctx, gconv.String(appId))

	if err != nil || app == nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "该商家应用不存在", dao.AlipayMerchantAppConfig.Table())
	}

	data := do.AlipayMerchantAppConfig{}

	if info.UserPolicy != "" {
		data.UserPolicy = info.UserPolicy
	}

	if info.PrivacyPolicy != "" {
		data.PrivacyPolicy = info.PrivacyPolicy
	}

	model := dao.AlipayMerchantAppConfig.Ctx(ctx)
	affected, err := daoctl.UpdateWithError(model.Data(data).OmitEmptyData().Where(do.AlipayMerchantAppConfig{AppId: appId}))

	if err != nil {
		return false, sys_service.SysLogs().ErrorSimple(ctx, err, "商家应用协议添加失败", dao.AlipayMerchantAppConfig.Table())
	}

	return affected > 0, nil
}

// GetPolicy 获取协议
func (s *sMerchantAppConfig) GetPolicy(ctx context.Context, appId string) (*alipay_model.GetPolicyRes, error) {
	res := alipay_model.GetPolicyRes{}

	err := dao.AlipayMerchantAppConfig.Ctx(ctx).Fields(dao.AlipayMerchantAppConfig.Columns().PrivacyPolicy, dao.AlipayMerchantAppConfig.Columns().UserPolicy).Where(do.AlipayMerchantAppConfig{
		AppId: appId,
	}).Scan(&res)

	if err != nil {
		return nil, sys_service.SysLogs().ErrorSimple(ctx, nil, "该AppId商家应用不存在", dao.AlipayMerchantAppConfig.Table())
	}

	return &res, nil
}
