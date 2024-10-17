package geoContainer

import (
	"context"
	"fmt"
	"net/http"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServerStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerStatusLogic {
	return &GetServerStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerStatusLogic) GetServerStatus(req *types.GetServerStatusReq) (*types.ErrorResponse, error) {
	geoServerURL := l.svcCtx.Config.GeoServerConfig.GeoServerURL
	username := l.svcCtx.Config.GeoServerConfig.Username
	password := l.svcCtx.Config.GeoServerConfig.Password
	workspace := l.svcCtx.Config.GeoServerConfig.Workspace
	mosaicName := "bev_" + req.TaskId
	storeName := mosaicName
	mosaicURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages/%s", geoServerURL, workspace, storeName, mosaicName)

	_req, err := http.NewRequest("GET", mosaicURL, nil)
	if err != nil {
		return nil, err
	}

	_req.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(_req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return &types.ErrorResponse{
			Code:    200,
			Message: "执行成功",
		}, nil

	case 500:
		return &types.ErrorResponse{
			Code:    500,
			Message: "数据源未找到",
		}, nil

	case 202:
		return &types.ErrorResponse{
			Code:    202,
			Message: "正在创建中",
		}, nil
	case 400:
		return &types.ErrorResponse{
			Code:    400,
			Message: "其他错误，查询日志",
		}, nil
	case 404:
		return &types.ErrorResponse{
			Code:    404,
			Message: "没有找到coverage",
		}, nil

	}

	return &types.ErrorResponse{
		Code:    400,
		Message: "其他错误，查询日志",
	}, nil
}
