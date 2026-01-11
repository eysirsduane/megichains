// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tron

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/backendadmin/internal/logic/tron"
	"megichains/apps/backendadmin/internal/svc"
)

func TronTransDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := tron.NewTronTransDetailLogic(r.Context(), svcCtx)
		err := l.TronTransDetail()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
