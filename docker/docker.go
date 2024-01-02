package docker

import (
	docker_utils "DistriAI-Node/docker/utils"
	"bufio"
	"context"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunScoreContainer() (float64, error) {

	image := "distrigroup/ml-device-score:v0.0.1"
	containerName := "ml-device-score"
	score := 0.0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return score, err
	}
	cli.NegotiateAPIVersion(ctx)

	containerID, err := docker_utils.RunContainer(ctx, cli, containerName,
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			// Easy debugging
			// 待定：在GPU上调试后再做决定
			// Resources: container.Resources{
			// 	DeviceRequests: []container.DeviceRequest{
			// 		{
			// 			Driver:       "", // or "nvidia"
			// 			Count:        -1,
			// 			Capabilities: [][]string{{"gpu"}},
			// 		},
			// 	},
			// },
			AutoRemove: true,
		})
	if err != nil {
		return score, err
	}

	reader, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true})
	if err != nil {
		return score, err
	}
	defer reader.Close()

	scanner1 := bufio.NewScanner(reader)
	for scanner1.Scan() {
		out := scanner1.Text()
		index := strings.Index(out, "Score:")
		if index > 0 {
			scoreStr := strings.TrimSpace(out[index+len("Score:"):])
			score, err = strconv.ParseFloat(scoreStr, 64)
			if err != nil {
				return score, err
			}
			return score, nil
		}
	}
	return score, nil
}