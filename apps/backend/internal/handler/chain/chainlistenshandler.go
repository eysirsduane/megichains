// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chain

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/backend/internal/logic/chain"
	"megichains/apps/backend/internal/svc"
)

func ChainListensHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := chain.NewChainListensLogic(r.Context(), svcCtx)
		resp, err := l.ChainListens()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
