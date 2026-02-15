// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chain

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/gateway/internal/logic/chain"
	"megichains/apps/gateway/internal/svc"
	"megichains/apps/gateway/internal/types"
)

func ChainListenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChainListenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chain.NewChainListenLogic(r.Context(), svcCtx)
		resp, err := l.ChainListen(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
