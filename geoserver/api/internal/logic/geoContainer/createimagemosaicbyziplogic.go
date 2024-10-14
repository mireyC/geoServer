package geoContainer

import (
	"context"
	"fmt"
	"geoserver/api/internal/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateImageMosaicByZipLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateImageMosaicByZipLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateImageMosaicByZipLogic {
	return &CreateImageMosaicByZipLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateImageMosaicByZip
// 根据上传的zip , 在指定的 store 下 创建一个ImageMosaic
// 参数：
//
//	storeName 仓库名称
//	filePath 上传zip文件路径
//
// 返回值:
//
//	success 是否成功
//	info    创建信息
func (l *CreateImageMosaicByZipLogic) CreateImageMosaicByZip(req *types.CreateImageMosaicByZipReq) (*types.ErrorResponse, error) {
	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	workSpace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeName := req.StoreName
	filePath := req.FilePath
	userName := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password

	//mosaicCheckURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages/%s.json")

	imageMosaicURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/file.imagemosaic?configure=all&coverageName=%s", geoServerURL, workSpace, storeName, storeName)
	fileData, err := os.Open(filePath)
	if err != nil {
		log.Println("open file err, ", err)
		return nil, err
	}

	client := &http.Client{}
	_req, err := http.NewRequest("PUT", imageMosaicURL, fileData)
	if err != nil {
		return nil, err
	}
	_req.Header.Set("Content-Type", "application/zip")
	_req.SetBasicAuth(userName, password)
	resp, err := client.Do(_req)
	if err != nil {
		return nil, err
	}

	// 201 (create) and 200 (ok)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(fmt.Errorf("failed to create/update mosaic, status: %s, body: %s", resp.Status, string(body)))
	}

	return util.ParseErrorCode(err), nil
}
