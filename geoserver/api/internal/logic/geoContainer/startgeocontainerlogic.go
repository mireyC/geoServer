package geoContainer

import (
	"context"
	"fmt"
	"geoserver/api/internal/svc"
	"geoserver/api/internal/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strings"
)

type StartGeoContainerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartGeoContainerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartGeoContainerLogic {
	return &StartGeoContainerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// StartGeoContainer
// 检查本地是否有 对应容器， 有则启动，
// 没有则 第一次请求 创建容器， 第二次请求 启动容器
func (l *StartGeoContainerLogic) StartGeoContainer() (*types.StartGeoContainerResp, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	imageName := l.svcCtx.Config.DockerImage.ImageName
	imageTag := l.svcCtx.Config.DockerImage.ImageTag
	port := l.svcCtx.Config.DockerImage.DockerHostPort
	source := l.svcCtx.Config.DockerImage.Source
	target := l.svcCtx.Config.DockerImage.Target

	// 检查是否已存在 oscarfonts/geoserver 的容器
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	fullImage := imageName + ":" + imageTag
	var containerID string
	for _, c := range containers {
		if c.Image == fullImage {
			containerID = c.ID
			break
		}
	}

	// 如果找到已有的容器，启动容器
	if containerID != "" {
		if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
			return nil, err
		}
		return &types.StartGeoContainerResp{
			Success: true,
			Info:    "contain start success",
		}, nil

	}

	// 如果没有找到容器，判断是否有镜像
	images, _ := cli.ImageList(ctx, image.ListOptions{All: true})

	found := false
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == fullImage {
				found = true
				break
			} else if strings.HasPrefix(tag, imageName) {
				msg := fmt.Sprintf("本地已存在镜像：%s, 与所需镜像：%s 不符, 请拉取所需镜像", tag, fullImage)
				return &types.StartGeoContainerResp{
					Success: false,
					Info:    msg,
				}, nil
			}
		}
		if found {
			break
		}
	}

	if !found {
		//go func() {
		//	reader, err := cli.ImagePull(ctx, "oscarfonts/geoserver", image.PullOptions{})
		//	if err != nil {
		//		log.Println("pull image: ", err)
		//	}
		//	defer reader.Close()
		//	io.Copy(os.Stdout, reader)
		//}()
		msg := fmt.Sprintf("本地不存在镜像：%s, 请拉取后重试", fullImage)
		return &types.StartGeoContainerResp{
			Success: false,
			Info:    msg,
		}, nil
	}

	// 创建容器
	// 映射端口 port，指定挂载配置

	go func() {
		hostBinding := nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: port,
		}
		portBinding := nat.PortMap{
			"8080/tcp": []nat.PortBinding{hostBinding},
		}

		mounts := []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: source, // 本地路径
				Target: target, // 容器内路径
			},
		}
		_, err := cli.ContainerCreate(ctx, &container.Config{
			Image: fullImage,
		}, &container.HostConfig{
			PortBindings: portBinding,
			Mounts:       mounts,
		}, nil, nil, "")
		if err != nil {
			log.Println("create err! ", err)
		}
	}()

	return &types.StartGeoContainerResp{
		Success: false,
		Info:    "创建容器中, 请稍后重试启动容器",
	}, nil
}
