package etc

import (
	"context"
	"fmt"
	"geoserver/api/internal/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"os"
)

func LoadImageAndCreateContainer(l config.Config) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	imageName := l.DockerImage.ImageName
	imageTag := l.DockerImage.ImageTag
	port := l.DockerImage.DockerHostPort
	source := l.DockerImage.Source
	target := l.DockerImage.Target
	fullImage := imageName + ":" + imageTag
	// 如果没有找到容器，判断是否有镜像
	images, _ := cli.ImageList(ctx, image.ListOptions{All: true})

	found := false
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == fullImage {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// 不存在拉取镜像
	if !found {
		fmt.Println("Image not found, loading from tar...")
		tarUrl := l.DockerImage.TarUrl
		tarFile, err := os.Open(tarUrl) // image tar 文件路径
		if err != nil {
			panic(err)
		}
		defer tarFile.Close()

		// 使用 Docker API 加载镜像

		loadResponse, err := cli.ImageLoad(ctx, tarFile, true)
		if err != nil {
			panic(err)
		}
		defer loadResponse.Body.Close()
		io.Copy(os.Stdout, loadResponse.Body) // 输出 Docker 加载日志到标准输出
		fmt.Println("Image loaded successfully")
	}

	// 检查是否已存在 oscarfonts/geoserver 的容器
	containers, _ := cli.ContainerList(ctx, container.ListOptions{All: true})
	var containerID string
	for _, c := range containers {
		if c.Image == fullImage {
			containerID = c.ID
			break
		}
	}

	// 不存在拉取容器
	if containerID == "" {
		containerName := l.DockerImage.ContainerName
		netWorkName := l.DockerImage.NetWorkName
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
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: fullImage,
		}, &container.HostConfig{
			PortBindings: portBinding,
			Mounts:       mounts,
			NetworkMode:  container.NetworkMode(netWorkName),
		}, nil, nil, containerName)
		if err != nil {
			log.Println("create err! ", err)
		}

		containerID = resp.ID // 更新容器 ID
		log.Printf("Created container %s\n", containerID)
	}

	err = cli.ContainerStart(ctx, containerID, container.StartOptions{})
	fmt.Println("err start: ", err)
}
