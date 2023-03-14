package boot

import (
    "context"
    "github.com/SupenBysz/gf-admin-community/sys_controller"
    "github.com/SupenBysz/gf-admin-community/sys_service"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/os/gcmd"
    "github.com/kysion/alipay-test/alipay_controller"
    merchant_controller "github.com/kysion/alipay-test/alipay_controller/merchant"
    _ "github.com/kysion/alipay-test/example/internal/boot/internal"
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
                    group.Group("/common/sys_file", func(group *ghttp.RouterGroup) { group.Bind(sys_controller.SysFile) })
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
                    alipay_controller.Gateway.AliPayServices,
                    alipay_controller.Gateway.AliPayCallback,
                    alipay_controller.Gateway.GetAlipayUserInfo,

                    // 异步通知
                    alipay_controller.MerchantNotify.NotifyServices,
                )

                // 直接通过回调获取用户信息
                //group.GET("/gateway.invite", func(r *ghttp.Request) {
                //	// 将URL拼接上unionMainID，然后进行URL编码  写一个conterrller  不建议在回调地址上面加自定义参数，实在需要可以写在state字段里面
                //	// alipays://platformapi/startapp?appId=2021003130652097&page=pages%2Fauthorize%2Findex%3FbizData%3D%7B%22platformCode%22%3A%22O%22%2C%22taskType%22%3A%22INTERFACE_AUTH%22%2C%22agentOpParam%22%3A%7B%22redirectUri%22%3A%22https%3A%2F%2Falipay.jditco.com%2Falipay%2Fgateway.callback%22%2C%22appTypes%22%3A%5B%22TINYAPP%22%2C%22WEBAPP%22%2C%22PUBLICAPP%22%2C%22MOBILEAPP%22%5D%2C%22isvAppId%22%3A%222021003179681073%22%7D%7D
                //	r.Response.RedirectTo("alipays://platformapi/startapp?appId=2021003130652097&page=pages/authorize/index?bizData=%7B%22platformCode%22%3A%22O%22%2C%22taskType%22%3A%22INTERFACE_AUTH%22%2C%22agentOpParam%22%3A%7B%22redirectUri%22%3A%22https%3A%2F%2Falipay.jditco.com%2Falipay%2Fgateway.callback%3FunionMainId%3D123123123%22%2C%22appTypes%22%3A%5B%22TINYAPP%22%2C%22WEBAPP%22%2C%22PUBLICAPP%22%2C%22MOBILEAPP%22%5D%2C%22isvAppId%22%3A%222021003179681073%22%7D%7D")
                //
                //})

                group.Group("/merchant", func(group *ghttp.RouterGroup) {
                    group.Bind(alipay_controller.AlipayMerchantAppConfig)

                })

                group.Group("/h5Pay", func(group *ghttp.RouterGroup) {
                    group.Bind(merchant_controller.AlipayMerchantH5Pay)
                })

                group.Group("/third", func(group *ghttp.RouterGroup) {
                    group.Bind(alipay_controller.AlipayThirdAppConfig)
                })

            })

            s.Run()
            return nil
        },
    }
)
