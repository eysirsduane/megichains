// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package solana

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/backendadmin/internal/logic/solana"
	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"
)

func SolanaTransListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SolanaTransListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := solana.NewSolanaTransListLogic(r.Context(), svcCtx)
		resp, err := l.SolanaTransList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
