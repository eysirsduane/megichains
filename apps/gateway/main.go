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

	"megichains/apps/gateway/internal/handler"
	"megichains/apps/gateway/internal/middleware"
	"megichains/apps/gateway/internal/svc"

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

	cfg.RestConf.Port = 7002
	server := rest.MustNewServer(cfg.RestConf)
	defer server.Stop()

	apimiidle := middleware.NewApiAccessPermissionMiddleware(db)

	userservice := service.NewUserService(db)
	excfgservice := service.NewRangeConfigService(db)
	evmservice := service.NewEvmService(db)
	addrservice := service.NewAddressService(db)
	orderservice := service.NewMerchOrderService(db)
	tronservice := service.NewTronService(db)
	chainservice := service.NewChainService(&cfg, db)
	solanaservice := service.NewSolanaService(db)
	listenservice := service.NewListenService(&cfg, db, addrservice, orderservice, chainservice, evmservice, tronservice, solanaservice)
	authservice := service.NewAuthService(db, cfg.Auth.AccessSecret, cfg.Auth.AccessExpire, cfg.Auth.RefreshSecret, cfg.Auth.RefreshExpire, cfg.Auth.Issuer)
	ctx := svc.NewServiceContext(cfg, apimiidle, excfgservice, userservice, authservice, addrservice, listenservice)
	handler.RegisterHandlers(server, ctx)

	httpx.SetOkHandler(biz.OkHandler)
	httpx.SetErrorHandlerCtx(biz.ErrHandler(cfg.Name))

	starting := fmt.Sprintf("Starting http server %s at %s:%d ...", cfg.Name, cfg.Host, cfg.Port)
	fmt.Println(starting)
	logx.Info(starting)

	server.Start()
}
