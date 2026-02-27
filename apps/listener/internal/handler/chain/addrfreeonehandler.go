// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chain

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/listener/internal/logic/chain"
	"megichains/apps/listener/internal/svc"
	"megichains/apps/listener/internal/types"
)

func AddrFreeoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddrFreeoneReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chain.NewAddrFreeoneLogic(r.Context(), svcCtx)
		resp, err := l.AddrFreeone(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
