package merchant

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/gopay"
)

// 转账
type sMerchantTransfer struct {
}

func init() {
	service.RegisterMerchantTransfer(NewMerchantTransfer())
}

func NewMerchantTransfer() *sMerchantTransfer {

	result := &sMerchantTransfer{}

	return result
}

// FundTransUniTransfer 单笔转账申请
func (s *sMerchantTransfer) FundTransUniTransfer(ctx context.Context, appId string, info *alipay_model.FundTransUniTransferReq) (aliRsp *alipay_model.TransUniTransferRes, err error) {
	// 商家AppId解析，获取商家应用，创建阿里支付客户端
	// appId, _ := strconv.ParseInt(g.RequestFromCtx(ctx).Get("appId").String(), 32, 0)

	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil {
		return nil, err
	}

	// 通过商家中的第三方应用的AppId创建客户端,但是AppAuthCode 需要是商家的
	client, err := aliyun.NewClient(ctx, merchantApp.ThirdAppId)

	// bindReq := alipay_model.TradeRelationBindReq{}
	// bindReq.OutRequestNo = gconv.String(orderId) // 外部请求号  == 单号

	data := kconv.Struct(info, &gopay.BodyMap{})
	data.Set("app_auth_token", merchantApp.AppAuthToken)

	res, err := client.FundTransUniTransfer(ctx, *data)

	ret := alipay_model.TransUniTransferRes{}
	gconv.Struct(&res.Response, &ret)

	return &ret, nil
}

// 余额查询
