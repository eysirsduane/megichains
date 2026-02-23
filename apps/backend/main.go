// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"megichains/pkg/biz"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service"

	"megichains/apps/backend/internal/handler"
	"megichains/apps/backend/internal/svc"

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
	server := rest.MustNewServer(cfg.RestConf, rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
		resp := biz.Response{
			Code: 401,
			Msg:  "登录状态失效，请重新登录",
			Data: nil,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Stop()

	excfgservice := service.NewRangeConfigService(db)
	authservice := service.NewAuthService(db, cfg.Auth.AccessSecret, cfg.Auth.AccessExpire, cfg.Auth.RefreshSecret, cfg.Auth.RefreshExpire, cfg.Auth.Issuer)
	userservice := service.NewUserService(db)
	merchservice := service.NewMerchService(db)
	addrservice := service.NewAddressService(db)
	orderservice := service.NewMerchOrderService(db)
	chainservice := service.NewChainService(&cfg, db)
	evmservice := service.NewEvmService(db)
	tronservice := service.NewTronService(db)
	solanaservice := service.NewSolanaService(db)
	listenservice := service.NewListenService(&cfg, db, merchservice, addrservice, orderservice, chainservice, evmservice, tronservice, solanaservice)
	fundservice := service.NewFundService(db)
	ctx := svc.NewServiceContext(cfg, excfgservice, authservice, userservice, merchservice, addrservice, orderservice, chainservice, listenservice, tronservice, evmservice, fundservice, solanaservice)
	handler.RegisterHandlers(server, ctx)

	httpx.SetOkHandler(biz.OkHandler)
	httpx.SetErrorHandlerCtx(biz.ErrHandler)

	starting := fmt.Sprintf("Starting http server %s at %s:%d ...", cfg.Name, cfg.Host, cfg.RestConf.Port)
	fmt.Println(starting)
	logx.Info(starting)

	server.Start()
}
