// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package fund

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/backendadmin/internal/logic/fund"
	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"
)

func AddressFundCollectLogDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddressFundCollectLogDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fund.NewAddressFundCollectLogDetailLogic(r.Context(), svcCtx)
		resp, err := l.AddressFundCollectLogDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
