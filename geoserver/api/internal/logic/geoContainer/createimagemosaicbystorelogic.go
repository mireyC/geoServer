package geoContainer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateImageMosaicByStoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateImageMosaicByStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateImageMosaicByStoreLogic {
	return &CreateImageMosaicByStoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateImageMosaicByStore
// 根据 storeName 创建 ImageMosaic
func (l *CreateImageMosaicByStoreLogic) CreateImageMosaicByStore(req *types.CreateImageMosaicByStoreReq) (*types.CreateImageMosaicByStoreResp, error) {
	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	workspace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeName := req.StoreName
	imageName := req.ImageName
	imageTitle := req.ImageTitle
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	coverageURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages", geoServerURL, workspace, storeName)

	coverageXML := fmt.Sprintf(`<coverage>
		<name>%s</name>
		<title>%s</title>
		<abstract>This is a sample image mosaic</abstract>
		<enabled>true</enabled>
		<store class="coverageStore">
			<name>%s</name>
		</store>
	</coverage>`, imageName, imageTitle, storeName)
	_req, err := http.NewRequest("POST", coverageURL, bytes.NewBufferString(coverageXML))
	if err != nil {
		return nil, err
	}

	_req.Header.Set("Content-Type", "text/xml")
	_req.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(_req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, err
	}

	return &types.CreateImageMosaicByStoreResp{
		Success: true,
		Info:    "create ImageMosaic Success",
	}, nil
}
