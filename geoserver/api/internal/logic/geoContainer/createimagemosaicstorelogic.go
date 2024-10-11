package geoContainer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateImageMosaicStoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateImageMosaicStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateImageMosaicStoreLogic {
	return &CreateImageMosaicStoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateImageMosaicStore
// 使用配置文件 的 username, password, workspace 及 用户传过来的 taskId bucketUrl 等信息 创建一个 store
func (l *CreateImageMosaicStoreLogic) CreateImageMosaicStore(req *types.CreateImageMosaicStoreReq) (*types.CreateImageMosaicStoreResp, error) {

	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	workSpace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeType := l.svcCtx.Config.GeoServerConfig.StoreType
	fileURL := req.BucketUrl

	coverageStoreName := "bev" + "_" + req.TaskId
	//fmt.Printf("coverageStoreName: ", coverageStoreName)
	// XML data for the coverage store
	coverageStoreXML := fmt.Sprintf(`<coverageStore>
  			<name>%s</name>
  			<type>%s</type>
  			<enabled>true</enabled>
  			<workspace>%s</workspace>
  			<url>%s</url>
	</coverageStore>`, coverageStoreName, storeType, workSpace, fileURL)

	// Create a new HTTP request
	createImageMosaicStoreURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores", geoServerURL, workSpace)
	_req, err := http.NewRequest("POST", createImageMosaicStoreURL, strings.NewReader(coverageStoreXML))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	// Set headers and authentication
	_req.Header.Set("Content-Type", "text/xml")
	_req.SetBasicAuth(username, password)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(_req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	return &types.CreateImageMosaicStoreResp{
		Success: true,
		Info:    "store create success",
	}, nil
}
