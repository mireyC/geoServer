package geoContainer

import (
	"context"
	"fmt"
	"geoserver/api/internal/util"
	"io/ioutil"
	"net/http"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelImageMosaicV2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelImageMosaicV2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *DelImageMosaicV2Logic {
	return &DelImageMosaicV2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DelImageMosaicV2
// 根据 taskId 删除 imageMosaic
func (l *DelImageMosaicV2Logic) DelImageMosaicV2(req *types.DelImageMosaicReqV2) (*types.ErrorResponse, error) {

	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	workspace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeName := "bev" + "_" + req.TaskId
	imageName := storeName
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	// the first del
	imageMosaicURL := fmt.Sprintf("%s//rest/workspaces/%s/layers/%s", geoServerURL, workspace, imageName)
	client := &http.Client{}
	delReq, err := http.NewRequest("DELETE", imageMosaicURL, nil)
	if err != nil {
		return nil, err
	}
	delReq.SetBasicAuth(username, password)
	resp, err := client.Do(delReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// the second del
	imageMosaicURL = fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages/%s",
		geoServerURL, workspace, storeName, imageName)

	client = &http.Client{}
	deleteReq, err := http.NewRequest("DELETE", imageMosaicURL, nil)
	if err != nil {
		return nil, err
	}

	deleteReq.SetBasicAuth(username, password)

	resp, err = client.Do(deleteReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the deletion was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {

		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		return nil, fmt.Errorf("failed to delete Image Mosaic, code: %d, resource not found", int(resp.StatusCode))
	}

	return util.ParseErrorCode(err), nil

}
