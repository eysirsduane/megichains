// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CheckMiddleware struct {
	db *gorm.DB
}

func NewCheckMiddleware(db *gorm.DB) *CheckMiddleware {
	return &CheckMiddleware{
		db: db,
	}
}

func (m *CheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		un := r.Context().Value("username").(string)
		logx.Infof("check middleware get username:%v", un)

		user := &entity.User{}
		if err := m.db.Where("username = ? and status = ?", un, global.UserStatusNormal).First(user).Error; err != nil {
			logx.Errorf("check middleware failed, username:%v, err:%v", un, err)

			resp := &CommonResp{Code: 1, Message: "Unauthorized"}

			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")

			w.Write(global.ObjToBytes(resp))

			return
		}

		next(w, r)
	}
}
