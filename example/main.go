package main

import (
	_ "github.com/SupenBysz/gf-admin-community"
	_ "github.com/SupenBysz/gf-admin-company-modules"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/kysion/alipay-library/example/internal/boot"
	_ "github.com/kysion/alipay-library/internal/logic"
	_ "github.com/kysion/base-library/base_hook"
)

func main() {
	boot.Main.Run(gctx.New())
}
