// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tron

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/backend/internal/logic/tron"
	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
)

func TronTransListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TronTransListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tron.NewTronTransListLogic(r.Context(), svcCtx)
		resp, err := l.TronTransList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
