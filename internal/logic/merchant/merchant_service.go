package merchant

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-test/alipay_model"
	"github.com/kysion/alipay-test/internal/logic/internal/aliyun"
	"time"
)

/*
	商户相关API服务
*/

type sMerchantService struct {
	redisCache *gcache.Cache
	Duration   time.Duration
}

func NewMerchantService() *sMerchantService {
	return &sMerchantService{
		redisCache: gcache.New(),
	}
}

func (s *sMerchantService) UserInfoAuth(ctx context.Context) (string, error) {
	// 后端需传递：商家参数、appID、
	// 前段提交数据：code、支付金额、相关信息

	// 等待前端提交code   --> 通过code拿到消费者token  --> 拼装支付相关参数  -->  调用支付接口
	client, err := aliyun.NewClient(ctx, "")

	bm, err := alipay.ParseNotifyToBodyMap(g.RequestFromCtx(ctx).Request)

	client.SetAppAuthToken(bm.Get("token"))

	// 还用不同的Clinet
	//response, err := client.UserInfoAuth(ctx, gopay.BodyMap{
	//	"scopes": []string{"auth_base"},
	//	"state":  "init",
	//})
	code := bm.Get("code")
	response, err := client.SystemOauthToken(ctx, gopay.BodyMap{
		"code":       code,
		"grant_type": "authorization_code",
	})
	fmt.Println(response)

	//
	//g.RequestFromCtx(ctx).Response.Header().Set("Content-Type", "text/html;charset=utf8")
	//g.RequestFromCtx(ctx).Response.WriteExit(response)

	return "", err
}

// TradeAppPay 测试手机APP支付
func (s *sMerchantService) TradeAppPay(ctx context.Context, info alipay_model.TradeAppPay) {
	client, _ := aliyun.NewClient(ctx, "")
	//请求参数
	bm := make(gopay.BodyMap)
	//bm.Set("subject", "测试APP支付")                     // 订单标题
	//bm.Set("out_trade_no", "GZ201901301040355706100469") // 商户网站唯一订单号
	//bm.Set("total_amount", "1.00")                       // 订单总金额，单位为元

	bm.Set("subject", info.Subject)                        // 订单标题
	bm.Set("out_trade_no", info.OutTradeNo)                // 商户网站唯一订单号
	bm.Set("total_amount", gconv.String(info.TotalAmount)) // 订单总金额，单位为元

	//手机APP支付参数请求
	payParam, err := client.TradeAppPay(ctx, bm)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("payParam:", payParam)
}

// TradeWapPay 手机网站支付
func (s *sMerchantService) TradeWapPay(ctx context.Context, info alipay_model.TradeWapPay) {
	//client, _ := NewClient(ctx)
	//
	//bm := make(gopay.BodyMap)
	//bm.Set("subject", "手机网站测试支付")
	//bm.Set("out_trade_no", "GZ201901301040355703")
	//bm.Set("quit_url", "https://alipay.jditco.com/gateway.callback")
	//bm.Set("total_amount", "100.00")
	//bm.Set("product_code", "QUICK_WAP_WAY")
	////手机网站支付请求
	//payUrl, err := client.TradeWapPay(ctx, bm)
	//if err != nil {
	//	xlog.Error("err:", err)
	//	return
	//}
	//xlog.Debug("payUrl:", payUrl)
}
