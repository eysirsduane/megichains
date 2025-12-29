// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"megichains/pkg/biz"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service"

	"megichains/apps/backendadmin/internal/handler"
	"megichains/apps/backendadmin/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "../../etc/megichains.dev.yaml", "the config file")

func main() {
	flag.Parse()

	var cfg global.BackendesConfig
	conf.MustLoad(*configFile, &cfg)
	logx.MustSetup(cfg.Log)
	defer logx.Close()

	db, err := entity.NewGormDB(&cfg)
	if err != nil {
		logx.Errorf("init gorm db failed, err:%v", err)
		panic(err)
	}

	cfg.RestConf.Port = 7001
	server := rest.MustNewServer(cfg.RestConf)
	defer server.Stop()

	excfgservice := service.NewRangeConfigService(db)
	authservice := service.NewAuthService(db, cfg.Auth.AccessSecret, cfg.Auth.AccessExpire, cfg.Auth.RefreshSecret, cfg.Auth.RefreshExpire, cfg.Auth.Issuer)
	userservice := service.NewUserService(db)
	orderservice := service.NewMerchOrderService(db)
	tronservice := service.NewTronService(db)
	evmservice := service.NewEvmService(db)
	ctx := svc.NewServiceContext(cfg, excfgservice, authservice, userservice, orderservice, tronservice, evmservice)
	handler.RegisterHandlers(server, ctx)

	httpx.SetOkHandler(biz.OkHandler)
	httpx.SetErrorHandlerCtx(biz.ErrHandler(cfg.Name))

	starting := fmt.Sprintf("Starting http server %s at %s:%d ...", cfg.Name, cfg.Host, cfg.RestConf.Port)
	fmt.Println(starting)
	logx.Info(starting)

	server.Start()
}
