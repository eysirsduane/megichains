// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package address

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megichains/apps/backend/internal/logic/address"
	"megichains/apps/backend/internal/svc"
)

func AddressGroupAllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewAddressGroupAllLogic(r.Context(), svcCtx)
		resp, err := l.AddressGroupAll()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
