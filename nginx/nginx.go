package nginx

import (
	"DistriAI-Node/utils/nginx_utils"
	"os/exec"
)

func StartNginx(nginxPort, consolePost, serverPost string) error {

	err := nginx_utils.GenNginxConfig(nginxPort, consolePost, serverPost)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "service", "nginx", "restart")
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func StopNginx() error {
	cmd := exec.Command("sudo", "service", "nginx", "stop")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}