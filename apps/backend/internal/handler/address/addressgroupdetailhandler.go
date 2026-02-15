// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package address

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/backend/internal/logic/address"
	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
)

func AddressGroupDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddressGroupDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := address.NewAddressGroupDetailLogic(r.Context(), svcCtx)
		resp, err := l.AddressGroupDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
