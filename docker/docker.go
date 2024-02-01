package docker

import (
	"DistriAI-Node/config"
	docker_utils "DistriAI-Node/docker/utils"
	"DistriAI-Node/pattern"
	logs "DistriAI-Node/utils/log_utils"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func RunScoreContainer(isGPU bool) (float64, error) {
	oldScore := 0.0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return oldScore, err
	}
	cli.NegotiateAPIVersion(ctx)

	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}
	if isGPU {
		hostConfig.Runtime = "nvidia"
		hostConfig.Resources = container.Resources{
			DeviceRequests: []container.DeviceRequest{
				{
					Count:        -1,
					Capabilities: [][]string{{"gpu"}},
				},
			},
		}
	}

	containerID, err := docker_utils.RunContainer(ctx, cli, pattern.SCORE_CONTAINER,
		&container.Config{
			Image: pattern.SCORE_NAME,
		},
		hostConfig)
	if err != nil {
		return oldScore, err
	}

	reader, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true})
	if err != nil {
		return oldScore, err
	}
	defer reader.Close()

	scanner1 := bufio.NewScanner(reader)
	for scanner1.Scan() {
		out := scanner1.Text()
		index := strings.Index(out, "Score:")
		if index > 0 {
			scoreStr := strings.TrimSpace(out[index+len("Score:"):])
			newScore, err := strconv.ParseFloat(scoreStr, 64)
			oldScore = (oldScore + newScore) / 2
			if err != nil {
				return oldScore, err
			}
		}
	}
	return oldScore, nil
}

func RunWorkspaceContainer(isGPU bool, mlToken string) (string, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	cli.NegotiateAPIVersion(ctx)

	containerName := pattern.ML_WORKSPACE_CONTAINER
	containerConfig := &container.Config{
		Image: pattern.ML_WORKSPACE_NAME,
		Env: []string{
			fmt.Sprintf("AUTHENTICATE_VIA_JUPYTER=%s", mlToken),
		},
		Tty:   true,
	}

	portBind := nat.PortMap{
		nat.Port("8080/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: config.GlobalConfig.Console.ConsolePost,
			},
		}}

	hostConfig := &container.HostConfig{
		PortBindings: portBind,
		Binds: []string{
			fmt.Sprintf("%s/ml-workspace:/workspace", config.GlobalConfig.Console.WorkDirectory),
			"myvolume:/data",
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		ShmSize: 512 * 1024 * 1024, // 512MB
	}
	if isGPU {
		containerName = pattern.ML_WORKSPACE_GPU_CONTAINER
		containerConfig.Image = pattern.ML_WORKSPACE_GPU_NAME
		hostConfig.Runtime = "nvidia"
		hostConfig.Resources = container.Resources{
			DeviceRequests: []container.DeviceRequest{
				{
					Count:        -1,
					Capabilities: [][]string{{"gpu"}},
				},
			},
		}
	}

	containerID, err := docker_utils.RunContainer(ctx, cli, containerName,
		containerConfig,
		hostConfig)
	if err != nil {
		return "", err
	}
	return containerID, nil
}

func StopWorkspaceContainer(containerID string) error {
	logs.Normal("Stop workspace container")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	if err := docker_utils.StopAndRemoveContainer(ctx, cli, containerID); err != nil {
		return err
	}

	dir := config.GlobalConfig.Console.WorkDirectory + "/ml-workspace"
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}
