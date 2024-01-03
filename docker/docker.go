package docker

import (
	docker_utils "DistriAI-Node/docker/utils"
	"DistriAI-Node/pattern"
	logs "DistriAI-Node/utils/log_utils"
	"bufio"
	"context"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunScoreContainer() (float64, error) {
	oldScore := 0.0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return oldScore, err
	}
	cli.NegotiateAPIVersion(ctx)

	containerID, err := docker_utils.RunContainer(ctx, cli, pattern.SCORE_CONTAINER,
		&container.Config{
			Image: pattern.SCORE_NAME,
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
		
		logs.Normal(out)

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
