package geoContainer

import (
	"context"
	"fmt"
	"net/http"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelImageMosaicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelImageMosaicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelImageMosaicLogic {
	return &DelImageMosaicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DelImageMosaic
// 删除指定 workspace 下的 coverage store，连带删除所有相关的 Image Mosaic 数据
// 参数:
//
//	storeName : coverage store name
//
// 返回值:
//
//	success
//	info
func (l *DelImageMosaicLogic) DelImageMosaic(req *types.DelImageMosaicReq) (*types.DelImageMosaicResp, error) {
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

	return &types.DelImageMosaicResp{
		Success: true,
		Info:    fmt.Sprintf("delete store %s and imageMosaic success", storeName),
	}, nil
}
