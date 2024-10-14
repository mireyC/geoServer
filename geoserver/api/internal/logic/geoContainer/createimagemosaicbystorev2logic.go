package geoContainer

import (
	"bytes"
	"context"
	"fmt"
	"geoserver/api/internal/util"
	"net/http"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateImageMosaicByStoreV2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateImageMosaicByStoreV2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateImageMosaicByStoreV2Logic {
	return &CreateImageMosaicByStoreV2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateImageMosaicByStoreV2
// 根据 taskID 创建 ImageMosaic
func (l *CreateImageMosaicByStoreV2Logic) CreateImageMosaicByStoreV2(req *types.CreateImageMosaicByStoreReqV2) (*types.ErrorResponse, error) {

	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	workspace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeName := "bev" + "_" + req.TaskId
	imageName := storeName
	imageTitle := storeName
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password

	// check if the imageMosaic already exists
	//imageCheckURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages/%s", geoServerURL, workspace, storeName, imageName)
	//client := &http.Client{}
	//checkReq, err := http.NewRequest("GET", imageCheckURL, nil)
	//if err != nil {
	//
	//	return nil, err
	//}
	//checkReq.SetBasicAuth(username, password)
	//checkResp, err := client.Do(checkReq)
	//if err != nil {
	//	return nil, err
	//}
	//defer checkResp.Body.Close()
	//if checkResp.StatusCode == http.StatusOK {
	//	return &types.ErrorResponse{
	//		Code:    500,
	//		Message: "Data Source already exists",
	//	}, nil
	//}

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

	if resp.StatusCode == http.StatusInternalServerError {
		return &types.ErrorResponse{Code: 500, Message: "Internal Server Error: " + "Data source already exists or store dose not exist"}, nil
	}

	return util.ParseErrorCode(err), nil
}
