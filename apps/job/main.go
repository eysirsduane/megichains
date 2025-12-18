package main

import (
	"flag"
	"fmt"
	"megichains/apps/job/keeps"
	"megichains/apps/job/manager"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "../../etc/megichains.dev.yaml", "the config file")

func main() {
	flag.Parse()

	var cfg global.Config
	conf.MustLoad(*configFile, &cfg)
	logx.MustSetup(cfg.Log)
	defer logx.Close()

	db, err := entity.NewGormDB(&cfg)
	if err != nil {
		logx.Errorf("init gorm db failed, err:%v", err)
		panic(err)
	}

	bsc := keeps.NewBSCMonitor(db)
	mgr := manager.NewKeepManager([]func(){bsc.Monitor})
	mgr.Start()

	starting := "Starting job..."
	fmt.Println(starting)
	logx.Info(starting)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	logx.Infof("收到信号:%s, 准备退出...", sig)
}
