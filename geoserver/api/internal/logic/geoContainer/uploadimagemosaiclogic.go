package geoContainer

import (
	"bytes"
	"context"
	"fmt"
	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadImageMosaicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImageMosaicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageMosaicLogic {
	return &UploadImageMosaicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadImageMosaic
// workspace 需要先创建好， 向指定 url 的store 上传 ImageMosaic, 如果 没有 该store会默认创建一个
func (l *UploadImageMosaicLogic) UploadImageMosaic(_req *types.UploadImageMosaicReq) (*types.UploadImageMosaicResp, error) {

	uploadZIPURL := _req.UploadZipUrl
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	zipFilePath := _req.ZipFilePath

	// Open the zip file
	file, err := os.Open(zipFilePath)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return nil, err
	}
	defer file.Close()

	// Read the file content
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading zip file:", err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", uploadZIPURL, bytes.NewBuffer(fileContent))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	// Set headers and authentication
	req.Header.Set("Content-Type", "application/zip")
	req.SetBasicAuth(username, password)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return &types.UploadImageMosaicResp{
		Success: true,
		Info:    "upload zip and create image success",
	}, nil
}
