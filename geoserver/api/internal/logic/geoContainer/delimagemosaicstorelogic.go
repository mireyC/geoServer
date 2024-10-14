package geoContainer

import (
	"context"
	"fmt"
	"geoserver/api/internal/util"
	"net/http"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelImageMosaicStoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelImageMosaicStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelImageMosaicStoreLogic {
	return &DelImageMosaicStoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DelImageMosaicStore
// 删除指定 workspace 下的 coverage store，连带删除所有相关的 Image Mosaic 数据
// 参数:
//
//	storeName : coverage store name
//
// 返回值:
//
//	success
//	info
func (l *DelImageMosaicStoreLogic) DelImageMosaicStore(req *types.DelImageMosaicStoreReq) (*types.ErrorResponse, error) {
	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	workspace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeName := req.StoreName
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	storeURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s?recurse=true", geoServerURL, workspace, storeName)

	client := &http.Client{}
	_req, err := http.NewRequest("DELETE", storeURL, nil)
	if err != nil {
		return nil, err
	}

	_req.SetBasicAuth(username, password)
	resp, err := client.Do(_req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return util.ParseErrorCode(err), nil

}
