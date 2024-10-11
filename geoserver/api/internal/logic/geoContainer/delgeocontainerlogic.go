package geoContainer

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log"

	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelGeoContainerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelGeoContainerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelGeoContainerLogic {
	return &DelGeoContainerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DelGeoContainer 删除所有 oscarfonts/geoserver 容器
func (l *DelGeoContainerLogic) DelGeoContainer() (*types.DelGeoContainerResp, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Fatalf("Error listing containers: %v", err)
	}

	go func() {
		imageName := l.svcCtx.Config.DockerImage.ImageName
		imageTag := l.svcCtx.Config.DockerImage.ImageTag
		fullImage := imageName + ":" + imageTag
		for _, containers := range containers {
			if containers.Image == fullImage {
				fmt.Printf("Deleting container %s...\n", containers.ID)
				err := cli.ContainerRemove(ctx, containers.ID, container.RemoveOptions{Force: true})
				if err != nil {
					log.Printf("Error removing container %s: %v\n", containers.ID, err)
				} else {
					fmt.Printf("Container %s deleted\n", containers.ID)
				}
			}
		}
	}()

	return &types.DelGeoContainerResp{
		Success: true,
		Info:    "container delete, please waiting",
	}, nil
}
