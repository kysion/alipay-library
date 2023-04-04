package boot

import (
	"context"
	"github.com/SupenBysz/gf-admin-community/sys_controller"
	"github.com/SupenBysz/gf-admin-community/sys_service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/kysion/alipay-library/alipay_controller"
	merchant_controller "github.com/kysion/alipay-library/alipay_controller/merchant"
	_ "github.com/kysion/alipay-library/example/internal/boot/internal"
)

// 需要让所有的配置从数据库加载

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			s.Group("/alipay", func(group *ghttp.RouterGroup) {
				// 注册中间件
				group.Middleware(
					sys_service.Middleware().CTX,
					sys_service.Middleware().ResponseHandler,
				)

				// 不需要鉴权，但是需要登录的路由
				group.Group("/", func(group *ghttp.RouterGroup) {
					// 注册中间件
					group.Middleware(
						sys_service.Middleware().Auth,
					)
					// 文件上传
					group.Group("/common/file", func(group *ghttp.RouterGroup) { group.Bind(sys_controller.SysFile) })
				})

				// 匿名路由绑定
				group.Group("/", func(group *ghttp.RouterGroup) {
					// 鉴权：登录，注册，找回密码等
					group.Group("/auth", func(group *ghttp.RouterGroup) { group.Bind(sys_controller.Auth) })
					// 图型验证码、短信验证码、地区
					group.Group("/common", func(group *ghttp.RouterGroup) {
						group.Bind(
							// 图型验证码
							sys_controller.Captcha,
							// 短信验证码
							sys_controller.SysSms,
							// 地区
							sys_controller.SysArea,
						)
					})
				})

				group.Bind(
					alipay_controller.Gateway.AliPayServices, // 消息接收
					alipay_controller.Gateway.AliPayCallback, // 网关回调

					// 异步通知
					alipay_controller.MerchantNotify.NotifyServices,

					// 商家授权 回调中接收
					merchant_controller.MerchantService.AuthMerchantApp,

					// 后端获取阿里用户信息 回调中接收
					merchant_controller.MerchantService.GetAlipayUserInfo,

					// 前端传递auth_code和appId获取用户信息
					//merchant_controller.MerchantService.GetUserInfoByAuthCode,
					
					// 版本提交
					merchant_controller.MerchantService.SubmitAppVersionAudit,

				)

				// 支付
				group.Group("/pay", func(group *ghttp.RouterGroup) {
					// h5支付
					group.Bind(merchant_controller.MerchantH5Pay)
					// 小程序支付
					//group.Bind(merchant_controller.MerchantH5Pay)

				})

				// 阿里商户应用配置表
				group.Group("/alipay_merchant", func(group *ghttp.RouterGroup) {
					group.Bind(alipay_controller.AlipayMerchantAppConfig)
				})

				// 阿里第三方应用配置表
				group.Group("/alipay_third", func(group *ghttp.RouterGroup) {
					group.Bind(alipay_controller.AlipayThirdAppConfig)
				})

				// 阿里消费者应用配置表
				group.Group("/alipay_consumer", func(group *ghttp.RouterGroup) {
				})

			})

			s.Run()
			return nil
		},
	}
)
