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

type ApiAccessPermissionMiddleware struct {
	db *gorm.DB
}

func NewApiAccessPermissionMiddleware(db *gorm.DB) *ApiAccessPermissionMiddleware {
	return &ApiAccessPermissionMiddleware{db: db}
}

func (m *ApiAccessPermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &CommonResp{Code: 1}

		maccount := r.Header.Get("merchant-account")
		csign := r.Header.Get("sign")
		if len(maccount) < 7 || len(csign) < 13 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "header keys not enough"

			w.Write(global.ObjToBytes(resp))

			return
		}

		merch := &entity.Merchant{}
		err := m.db.Model(&entity.Merchant{}).Where("merchant_account = ?", maccount).First(merch).Error
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "merchant not found"

			w.Write(global.ObjToBytes(resp))

			return
		}

		bbytes, err := io.ReadAll(r.Body)
		if err != nil {
			logx.Errorf("api access permission middleware read body failed, err:%v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "request body invalid"

			w.Write(global.ObjToBytes(resp))
			return
		}

		logx.Infof("api access permission check request, mchaccount:%v, header:%v, body:%v, csign:%v", maccount, r.Header, string(bbytes), csign)

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
			logx.Errorf("api access permission middleware parse request body failed, err:%v", err)

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
			logx.Errorf("order create failed, mono:%v, receiver:%v, err:%v", req.MerchantOrderNo, req.Receiver, err)
			return
		}

		reqmap := make(map[string]any)
		reqmap["url"] = r.URL.String()
		reqmap["method"] = r.Method
		reqmap["ip"] = r.RemoteAddr
		reqmap["time"] = time.Now().Format("2006-01-02 15:04:05")
		reqmap["header"] = r.Header
		reqmap["query"] = r.URL.Query()
		reqmap["body"] = string(bbytes)

		r.Body = io.NopCloser(bytes.NewBuffer(bbytes))

		recorder := &ResponseRecorder{ResponseWriter: w}

		next(recorder, r)

		rlog := &entity.MerchantOrderRequestLog{
			MerchantOrderId: order.Id,
			Request:         global.ObjToJsonString(reqmap),
			Response:        recorder.body.String(),
			Description:     "",
		}

		logx.Infof("api access permission check response, mchaccount:%v, header:%+v, body:%v", maccount, w.Header(), recorder.body.String())

		err = m.db.Create(rlog).Error
		if err != nil {
			logx.Errorf("api access permission middleware create request log failed, err:%v", err)
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
