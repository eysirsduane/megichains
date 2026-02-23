// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/crypt"
	"megichains/pkg/global"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OrderTestPlaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderTestPlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderTestPlaceLogic {
	return &OrderTestPlaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderTestPlaceLogic) OrderTestPlace(req *types.OrderTestPlaceReq) (resp *types.Response, err error) {
	reqc := &converter.ChainListenReq{}
	copier.Copy(reqc, req)

	merch, err := l.svcCtx.MerchService.GetByAccount(reqc.MerchantAccount)
	if err != nil {
		return
	}

	reqc.MerchantOrderNo = fmt.Sprintf("TEST%v", time.Now().UnixMicro())

	rbytes := global.ObjToBytes(reqc)

	request, err := http.NewRequest(http.MethodPost, l.svcCtx.Config.Gateway.ListenApi, bytes.NewReader(rbytes))
	if err != nil {
		return
	}

	sign := crypt.HmacSHA256(string(rbytes), merch.SecretKey)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Merchant-Account", merch.MerchantAccount)
	request.Header.Add("Sign", sign)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	bbytes, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		logx.Errorf("merchant order place response is not ok, status:%v", response.StatusCode)
		err = biz.OrderTestPlaceFailed
		return
	}

	logx.Infof("merchant order place response body:%v", string(bbytes))

	resp = &types.Response{}
	global.BytesToObj(bbytes, resp)

	if resp.Code != 0 {
		logx.Errorf("merchant order place failed, code:%v, msg:%v", resp.Code, resp.Msg)
		err = biz.NewSpecificError(int64(resp.Code), resp.Msg)
		return
	}

	return
}
