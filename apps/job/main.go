package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"megichains/apps/job/keeps"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "../../etc/megichains.dev.yaml", "the config file")
var monitor *keeps.ChainMonitor

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

	ethservice := service.NewEvmService(db)
	solaservice := service.NewSolanaService(db)
	addrservice := service.NewAddressService(db)
	monitor = keeps.NewChainMonitor(&cfg, ethservice, addrservice, solaservice)

	starting := "Starting job..."
	fmt.Println(starting)
	logx.Info(starting)

	http.HandleFunc("/listen", listen)
	http.HandleFunc("/listens", listens)
	fmt.Println("HTTP 服务启动 :7002...")
	logx.Infof("HTTP 服务启动 :7002...")
	err = http.ListenAndServe(":7002", nil)
	if err != nil {
		fmt.Println("Error starting http server:", err)
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	logx.Infof("收到信号:%s, 准备退出...", sig)
}
func listens(w http.ResponseWriter, r *http.Request) {
	monitor.RangeListen()

	fmt.Fprintf(w, "已启动批量监听")
}

func listen(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	resp := &ResponseMessage{
		Code:    500,
		Message: "",
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		resp.Message = "读取请求体失败"
		returnError(w, resp)
		return
	}

	var req global.ListenReq
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		resp.Message = "解析请求体失败"
		returnError(w, resp)
		return
	}
	if req.MerchOrderId == "" {
		resp.Message = "订单号不能为空"
		returnError(w, resp)
		return
	}
	if req.Chain == "" {
		resp.Message = "链不能为空"
		returnError(w, resp)
		return
	}
	if req.Currency == "" {
		resp.Message = "币种不能为空"
		returnError(w, resp)
		return
	}
	if req.Seconds > 3600 || req.Seconds <= 0 {
		resp.Message = "监听时间不合法"
		returnError(w, resp)
		return
	}
	if req.Receiver == "" {
		resp.Message = "接收地址不能为空"
		returnError(w, resp)
		return
	}

	exist := startListen(global.ChainName(req.Chain), &req)
	if exist {
		logx.Errorf("监听地址已存在, receiver:%s", req.Receiver)
		resp.Message = "监听地址已存在"
		returnError(w, resp)
		return
	}

	resp.Code = 0
	resp.Message = fmt.Sprintf("监听启动成功, chain:%v, receiver:%s", req.Chain, req.Receiver)

	byts, _ := json.Marshal(resp)
	fmt.Fprint(w, string(byts))
}

func returnError(w http.ResponseWriter, resp *ResponseMessage) {
	byts, err := json.Marshal(resp)
	if err != nil {
		return
	}

	http.Error(w, string(byts), http.StatusInternalServerError)
}

func startListen(chain global.ChainName, req *global.ListenReq) (exist bool) {
	rkey := global.GetOrderAddressKey(req.Receiver, req.Currency)
	monitor.Receivers.Range(func(key, val any) bool {
		if key == rkey {
			logx.Infof("已存在监听地址, chain:%v, receiver:%v, currency:%v", chain, req.Receiver, req.Currency)
			exist = true
			return false
		}

		return true
	})

	if exist {
		return
	}

	go monitor.Listen(chain, req.Currency, req.MerchOrderId, req.Receiver, req.Seconds+120)

	return
}

type ResponseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
