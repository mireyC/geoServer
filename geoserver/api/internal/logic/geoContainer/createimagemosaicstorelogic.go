package geoContainer

import (
	"bytes"
	"context"
	"fmt"
	"geoserver/api/internal/util"
	"log"
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
func (l *CreateImageMosaicStoreLogic) CreateImageMosaicStore(req *types.CreateImageMosaicStoreReq) (*types.ErrorResponse, error) {

	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	workSpace := l.svcCtx.Config.GeoServerConfig.Workspace
	storeType := l.svcCtx.Config.GeoServerConfig.StoreType
	fileURL := req.BucketUrl

	coverageStoreName := "bev" + "_" + req.TaskId
	// check if store already exists
	checkStoreURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s", geoServerURL, workSpace, coverageStoreName)
	checkReq, err := http.NewRequest("GET", checkStoreURL, nil)
	if err != nil {
		fmt.Println("err check store, ", err)
		return nil, err
	}

	checkReq.SetBasicAuth(username, password)
	client := &http.Client{}
	checkResp, err := client.Do(checkReq)
	if err != nil {
		fmt.Println("Error checking coverage store: ", err)
		return nil, err
	}

	defer checkResp.Body.Close()

	if checkResp.StatusCode == http.StatusOK {
		return &types.ErrorResponse{
			Code:    500,
			Message: "Data Source already exists",
		}, nil
	}

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
	client = &http.Client{}
	resp, err := client.Do(_req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 调用创建服务， 创建 ImageMosaic
	go func() {

		imageName := coverageStoreName
		imageTitle := coverageStoreName
		storeName := coverageStoreName

		imageMosaicURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages", geoServerURL, workSpace, storeName)
		imageMosaicXML := fmt.Sprintf(`<coverage>
			<name>%s</name>
			<title>%s</title>
			<abstract>This is a sample image mosaic</abstract>
			<enabled>true</enabled>
			<store class="coverageStore">
				<name>%s</name>
			</store>
		</coverage>`, imageName, imageTitle, storeName)
		_req, er := http.NewRequest("POST", imageMosaicURL, bytes.NewBufferString(imageMosaicXML))
		if er != nil {
			log.Fatal("err create imageMosaic request during create store: , ", er)
		}

		_req.Header.Set("Content-Type", "text/xml")
		_req.SetBasicAuth(username, password)
		client = &http.Client{}
		resp, er := client.Do(_req)
		if er != nil {
			log.Fatal("err create imageMosaic during create store: , ", er)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {

			log.Fatal("err create imageMosaic during create store: , ", er)

		}

	}()

	return util.ParseErrorCode(err), nil
}
