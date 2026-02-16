// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"bytes"
	"io"
	"megichains/pkg/crypt"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"net/http"

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

		logx.Infof("api access permission check, mchaccount:%v, request body:%v, csign:%v", maccount, string(bbytes), csign)

		ok := crypt.VerifySignature(csign, bbytes, merch.SecretKey)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "signature invalid"

			w.Write(global.ObjToBytes(resp))

			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bbytes))

		next(w, r)
	}
}
