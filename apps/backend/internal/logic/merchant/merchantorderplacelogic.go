// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"
	"megichains/pkg/crypt"
	"megichains/pkg/global"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantOrderPlaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantOrderPlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantOrderPlaceLogic {
	return &MerchantOrderPlaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantOrderPlaceLogic) MerchantOrderPlace(req *types.MerchantOrderPlaceReq) (resp *types.Response, err error) {
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
		return
	}

	logx.Infof("merchant order place response body:%v", string(bbytes))

	return
}
