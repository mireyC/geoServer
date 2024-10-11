package geoContainer

import (
	"net/http"

	"geoserver/api/internal/logic/geoContainer"
	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateImageMosaicByStoreHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateImageMosaicByStoreReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := geoContainer.NewCreateImageMosaicByStoreLogic(r.Context(), svcCtx)
		resp, err := l.CreateImageMosaicByStore(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
