package gateway

import (
    "context"
    "github.com/SupenBysz/gf-admin-community/sys_service"
    "github.com/gogf/gf/v2/os/gcache"
    "github.com/gogf/gf/v2/util/gconv"
    "github.com/kysion/alipay-test/alipay_model"
    dao "github.com/kysion/alipay-test/alipay_model/alipay_dao"
    do "github.com/kysion/alipay-test/alipay_model/alipay_do"
    entity "github.com/kysion/alipay-test/alipay_model/alipay_entity"
    "github.com/kysion/base-library/utility/daoctl"
    "github.com/yitter/idgenerator-go/idgen"
    "time"
)

type sConsumer struct {
    redisCache *gcache.Cache
    Duration   time.Duration
}

func NewConsumerConfig() *sConsumer {
    return &sConsumer{
        redisCache: gcache.New(),
    }
}

// GetConsumerById 根据id查找消费者信息
func (s *sConsumer) GetConsumerById(ctx context.Context, id int64) (*alipay_model.AlipayConsumerConfig, error) {
    return daoctl.GetByIdWithError[alipay_model.AlipayConsumerConfig](dao.AlipayConsumerConfig.Ctx(ctx), id)
}

// GetConsumerBySysUserId  根据用户id查询消费者信息
func (s *sConsumer) GetConsumerBySysUserId(ctx context.Context, sysUserId int64) (*alipay_model.AlipayConsumerConfig, error) {
    result := alipay_model.AlipayConsumerConfig{}
    model := dao.AlipayConsumerConfig.Ctx(ctx)

    err := model.Where(dao.AlipayConsumerConfig.Columns().SysUserId, sysUserId).Scan(&result)

    if err != nil {
        return nil, err
    }

    return &result, nil
}

// CreateConsumer  创建消费者信息
func (s *sConsumer) CreateConsumer(ctx context.Context, info alipay_model.AlipayConsumerConfig) (*alipay_model.AlipayConsumerConfig, error) {
    data := do.AlipayConsumerConfig{}

    gconv.Struct(info, &data)

    data.Id = idgen.NextId()
    data.UserState = 1 // 用户状态默认正常

    if info.ExtJson == "" {
        data.ExtJson = nil
    }
    affected, err := daoctl.InsertWithError(
        dao.AlipayConsumerConfig.Ctx(ctx),
        data,
    )

    if affected == 0 || err != nil {
        return nil, sys_service.SysLogs().ErrorSimple(ctx, err, "消费者信息创建失败", dao.AlipayConsumerConfig.Table())
    }

    return s.GetConsumerById(ctx, gconv.Int64(data.Id))
}

// UpdateConsumer 更新消费者信息
func (s *sConsumer) UpdateConsumer(ctx context.Context, id int64, info alipay_model.UpdateConsumerReq) (bool, error) {
    // 首先判断消费者信息是否存在
    consumerInfo, err := daoctl.GetByIdWithError[entity.AlipayConsumerConfig](dao.AlipayConsumerConfig.Ctx(ctx), id)
    if err != nil || consumerInfo == nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "该消费者不存在", dao.AlipayConsumerConfig.Table())
    }
    data := do.AlipayConsumerConfig{}
    gconv.Struct(info, &data)

    model := dao.AlipayConsumerConfig.Ctx(ctx)
    affected, err := daoctl.UpdateWithError(model.Data(model).OmitNilData().Where(do.AlipayConsumerConfig{Id: id}))

    if err != nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "消费者信息更新失败", dao.AlipayConsumerConfig.Table())
    }

    return affected > 0, nil
}

// UpdateConsumerState 修改用户状态
func (s *sConsumer) UpdateConsumerState(ctx context.Context, id int64, state int) (bool, error) {
    affected, err := daoctl.UpdateWithError(dao.AlipayConsumerConfig.Ctx(ctx).Data(do.AlipayConsumerConfig{
        UserState: state,
    }).OmitNilData().Where(do.AlipayConsumerConfig{Id: id}))

    if err != nil {
        return false, sys_service.SysLogs().ErrorSimple(ctx, err, "消费者状态修改失败", dao.AlipayConsumerConfig.Table())
    }
    return affected > 0, err
}
