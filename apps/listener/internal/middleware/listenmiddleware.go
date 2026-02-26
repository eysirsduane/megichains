// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"bytes"
	"io"
	"megichains/pkg/converter"
	"megichains/pkg/crypt"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"net/http"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ListenMiddleware struct {
	db *gorm.DB
}

func NewListenMiddleware(db *gorm.DB) *ListenMiddleware {
	return &ListenMiddleware{db: db}
}

func (m *ListenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &CommonResp{Code: 1}

		maccount := r.Header.Get("merchant-account")
		csign := r.Header.Get("sign")
		if len(maccount) < 7 || len(csign) < 13 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "header invalid"

			w.Write(global.ObjToBytes(resp))

			return
		}

		merch := &entity.Merchant{}
		err := m.db.Model(&entity.Merchant{}).Where("merchant_account = ?", maccount).First(merch).Error
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "merchant invalid"

			w.Write(global.ObjToBytes(resp))

			return
		}

		bbytes, err := io.ReadAll(r.Body)
		if err != nil {
			logx.Errorf("listen middleware read body failed, err:%v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "request body invalid"

			w.Write(global.ObjToBytes(resp))
			return
		}

		logx.Infof("listen middleware check request, mchaccount:%v, header:%v, body:%v, csign:%v", maccount, r.Header, string(bbytes), csign)

		ok := crypt.VerifySignature(csign, bbytes, merch.SecretKey)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "signature invalid"

			w.Write(global.ObjToBytes(resp))

			return
		}

		req := &converter.ChainListenReq{}
		err = global.BytesToObj(bbytes, req)
		if err != nil {
			logx.Errorf("listen middleware parse request body failed, err:%v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "request body invalid"

			w.Write(global.ObjToBytes(resp))
			return
		}

		node, _ := snowflake.NewNode(1)
		sid := node.Generate()

		if req.Mode != string(global.OrderModeTest) {
			req.Mode = string(global.OrderModeProd)
		}
		order := &entity.MerchantOrder{
			OrderNo:         sid.String(),
			MerchantAccount: maccount,
			MerchantOrderNo: req.MerchantOrderNo,
			Chain:           string(req.Chain),
			Typo:            string(global.OrderTypoIn),
			Mode:            req.Mode,
			Status:          string(global.OrderStatusCreated),
			NotifyStatus:    string(global.NotifyStatusUnknown),
			Currency:        string(req.Currency),
			ToAddress:       req.Receiver,
			Description:     "",
		}
		err = m.db.Create(order).Error
		if err != nil {
			logx.Errorf("listen middleware order create failed, mono:%v, receiver:%v, err:%v", req.MerchantOrderNo, req.Receiver, err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "order invalid"

			w.Write(global.ObjToBytes(resp))

			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bbytes))
		recorder := &ResponseRecorder{ResponseWriter: w}

		reqmap := make(map[string]any)
		reqmap["url"] = r.URL.String()
		reqmap["method"] = r.Method
		reqmap["ip"] = r.RemoteAddr
		reqmap["time"] = time.Now().Format(global.DateTimeFormat)
		reqmap["header"] = r.Header
		reqmap["query"] = r.URL.Query()
		reqmap["body"] = string(bbytes)

		ilog := &entity.MerchantOrderInteractionLog{
			MerchantOrderId:       order.Id,
			PlaceRequest:          global.ObjToJsonString(reqmap),
			PlaceRequestTimestamp: uint64(time.Now().UnixMilli()),

			Description: "",
		}

		next(recorder, r)

		logx.Infof("listen middleware check response, mchaccount:%v, header:%+v, body:%v", maccount, w.Header(), recorder.body.String())

		ilog.PlaceResponse = recorder.body.String()
		ilog.PlaceResponseTimestamp = uint64(time.Now().UnixMilli())

		err = m.db.Create(ilog).Error
		if err != nil {
			logx.Errorf("listen middleware create interaction log failed, err:%v", err)
			return
		}
	}
}

type ResponseRecorder struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
