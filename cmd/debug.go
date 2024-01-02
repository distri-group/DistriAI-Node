package cmd

import (
	"DistriAI-Node/docker"
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"github.com/urfave/cli"
)

var DebugCommand = cli.Command{
	Name:  "debug",
	Usage: "debug.",
	Action: func(c *cli.Context) error {

		score, err := docker.RunScoreContainer()
		if err != nil {
			return err
		}
		logs.Normal(fmt.Sprintf("debug score: %v", score))
		return nil
	},
}
