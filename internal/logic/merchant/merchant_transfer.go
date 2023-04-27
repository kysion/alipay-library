package merchant

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/base-library/utility/kconv"
	"github.com/kysion/gopay"
	"github.com/kysion/gopay/pkg/xlog"
)

// 转账  仅支持商家自调用
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
	sys_service.SysLogs().InfoSimple(ctx, nil, "\n-------单笔转账申请 ------- ", "sMerchantTransfer")

	// 商家AppId解析，获取商家应用，创建阿里支付客户端
	// appId, _ := strconv.ParseInt(g.RequestFromCtx(ctx).Get("appId").String(), 32, 0)

	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil {
		return nil, err
	}

	// 只能通过商家自己创建client，然后
	client, err := aliyun.NewMerchantClient(ctx, merchantApp.AppId)

	reqData := alipay_model.FundTransUniTransferReq{
		OutBizNo:    info.OutBizNo,          // 商家侧唯一订单号，商家自定义
		TransAmount: info.TransAmount,       // 订单总金额，单位元
		ProductCode: "TRANS_ACCOUNT_NO_PWD", // 销售产品码
		BizScene:    "DIRECT_TRANSFER",      // 业务场景
		OrderTitle:  info.OrderTitle,        // 转账标题
		PayeeInfo: alipay_model.PayeeInfo{ // 收款方信息
			Identity:     info.PayeeInfo.Identity, // 参与方的标识ID，例如收款账号
			IdentityType: "ALIPAY_LOGON_ID",       // 参与方的标识类型
			Name:         info.PayeeInfo.Name,     // 参与方真实姓名
		}, // 收款方信息
		Remark:         info.Remark,
		BusinessParams: "",
	}

	data := kconv.Struct(reqData, &gopay.BodyMap{})
	data.Set("app_auth_token", merchantApp.AppAuthToken)

	res, err := client.FundTransUniTransfer(ctx, *data)

	ret := alipay_model.TransUniTransferRes{}
	gconv.Struct(&res.Response, &ret)

	return &ret, nil
}

// FundTransCommonQuery 查询转账详情  orderId != OutBizNo  调用转账成功后会返回order_id
func (s *sMerchantTransfer) FundTransCommonQuery(ctx context.Context, appId string, outBizNo string) (aliRsp *alipay_model.FundTransCommonQueryRes, err error) {
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil {
		return nil, err
	}

	// 只能通过商家自己创建client，然后
	client, err := aliyun.NewMerchantClient(ctx, merchantApp.AppId)

	bm := make(gopay.BodyMap)
	bm.Set("product_code", "TRANS_ACCOUNT_NO_PWD").
		Set("biz_scene", "DIRECT_TRANSFER").
		Set("out_biz_no", outBizNo)
	//Set("order_id", orderId)

	res, err := client.FundTransCommonQuery(ctx, bm)
	if err != nil {
		xlog.Error(err)
		return
	}

	ret := alipay_model.FundTransCommonQueryRes{}
	gconv.Struct(&res.Response, &ret)

	xlog.Debug("aliRsp:", *aliRsp)
	xlog.Debug("aliRsp.Response:", aliRsp.Response)

	return &ret, nil
}

// FundAccountQuery 余额查询
func (s *sMerchantTransfer) FundAccountQuery(ctx context.Context, appId string, userId string) (aliRsp *alipay_model.FundAccountQueryResponse, err error) {
	merchantApp, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, appId)
	if err != nil {
		return nil, err
	}

	// 只能通过商家自己创建client，然后
	client, err := aliyun.NewMerchantClient(ctx, merchantApp.AppId)

	bm := make(gopay.BodyMap)
	bm.Set("alipay_user_id", userId) /*.Set("account_type", "ACCTRANS_ACCOUNT")*/

	res, err := client.FundAccountQuery(ctx, bm)
	if err != nil {
		xlog.Error(err)
		return
	}

	ret := alipay_model.FundAccountQueryResponse{}
	gconv.Struct(&res.Response, &ret)

	xlog.Debug("aliRsp:", *aliRsp)
	xlog.Debug("aliRsp.Response:", aliRsp.Response)

	return &ret, nil
}
