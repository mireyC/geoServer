package geoContainer

import (
	"net/http"

	"geoserver/api/internal/logic/geoContainer"
	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadImageMosaicHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadImageMosaicReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := geoContainer.NewUploadImageMosaicLogic(r.Context(), svcCtx)
		resp, err := l.UploadImageMosaic(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
