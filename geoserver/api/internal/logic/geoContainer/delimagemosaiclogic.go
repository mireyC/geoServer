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
// 根据 StoreName, ImageName 删除 imageMosaic
func (l *DelImageMosaicLogic) DelImageMosaic(req *types.DelImageMosaicReq) (*types.ErrorResponse, error) {

	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	workspace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeName := req.StoreName
	imageName := req.ImageName
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password

	// first del
	imageMosaicURL := fmt.Sprintf("%s/rest/workspaces/%s/layers/%s", geoServerURL, workspace, imageName)

	client := &http.Client{}
	delRep, err := http.NewRequest("DELETE", imageMosaicURL, nil)
	if err != nil {
		return nil, err
	}
	delRep.SetBasicAuth(username, password)
	resp, err := client.Do(delRep)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Construct the URL for the DELETE request to delete the Image Mosaic
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
		return nil, fmt.Errorf("failed to delete Image Mosaic, code: %d,resource not found ", resp.StatusCode)
	}

	return util.ParseErrorCode(err), nil

}
