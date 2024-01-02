package docker_utils

import (
	logs "DistriAI-Node/utils/log_utils"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// 查询镜像是否已存在
func ImageExist(ctx context.Context, cli *client.Client, imageName string) (bool, string) {
	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return false, ""
	}

	for _, image := range images {
		for _, name := range image.RepoTags {
			if strings.Contains(imageName, name) {
				logs.Normal(fmt.Sprintf("Image %s exists", imageName))
				return true, image.ID
			}
		}
	}
	logs.Normal(fmt.Sprintf("Image %s does not exist", imageName))
	return false, ""
}

// 拉取镜像
func PullImage(ctx context.Context, cli *client.Client, imageName string) error {
	cmd := exec.Command("docker", "pull", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error : Start pulling image: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error : Wait pulling image: %v", err)
	}
	return nil
}

// 创建容器
func CreateContainer(ctx context.Context, cli *client.Client, containerName string, config *container.Config, hostConfig *container.HostConfig) (string, error) {
	resp, err := cli.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// 运行容器
func RunContainer(ctx context.Context, cli *client.Client, containerName string, config *container.Config, hostConfig *container.HostConfig) (string, error) {
	imageName := config.Image
	isCreated, _ := ImageExist(ctx, cli, imageName)
	if !isCreated {
		if err := PullImage(ctx, cli, imageName); err != nil {
			return "", err
		}
	}

	hostConfig.AutoRemove = true
	containerID, err := CreateContainer(ctx, cli, containerName, config, hostConfig)
	if err != nil {
		return "", err
	}

	logs.Normal(fmt.Sprintf("Start running container %s", containerName))

	if err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return containerID, nil
}

func DeleteImage(ctx context.Context, cli *client.Client, imageName string) error {
	isCreated, imageId := ImageExist(ctx, cli, imageName)
	if isCreated {
		_, err := cli.ImageRemove(ctx, imageId, types.ImageRemoveOptions{Force: true})
		if err != nil {
			return err
		}
	}
	return nil
}
