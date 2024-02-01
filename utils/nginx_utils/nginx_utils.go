package nginx_utils

import (
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"os"
)

func GenNginxConfig(nginxPort, consolePort, serverPost string) error {
	logs.Normal(fmt.Sprintf("Start nginx config, nginxPort: %v, consolePort: %v, serverPost: %v", nginxPort, consolePort, serverPost))
	dir := "/etc/nginx/sites-enabled"

	files, err := os.ReadDir(dir)
	if err != nil {
		logs.Error(fmt.Sprintf("Error reading config file: %v", err))
		return err
	}

	for _, file := range files {
		err = os.Remove(dir + "/" + file.Name())
		if err != nil {
			logs.Error(fmt.Sprintf("Error reading config file: %v", err))
			return err
		}
	}

	nginxConfig := fmt.Sprintf(`server {
	listen %v;
	listen [::]:%v;

	server_name distri-ai-node;

	location ^~ /distri/ {
		proxy_pass http://127.0.0.1:%v/;
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
		}

	location / {
		proxy_pass http://127.0.0.1:%v;
		proxy_http_version 1.1;
		proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}`, nginxPort, nginxPort, serverPost, consolePort)

	err = os.WriteFile("/etc/nginx/sites-enabled/distri", []byte(nginxConfig), 0644)
	if err != nil {
		logs.Error(fmt.Sprintf("Error reading config file: %v", err))
		return err
	}

	return nil
}
