# DistriAI-Node
Share your unused computing capacity to provide support for more AI creators in need and earn profits at the same time.

## Recommended configuration
### Software
- Operating System: Ubuntu 22.04
- Go: 1.21.0
### Hardware
- CPU: 4 cores
- Memory: 8GB
- Disk: 100GB
- GPU: Nvidia RTX 3090
- Network: 1Gbps

## Necessary software packages
Please ensure that the following software packages are installed before use.

- curl
- speedtest-cli
- docker
- nginx

***During the DistriAI-node program execution, all Docker containers will be stopped and the nginx configuration file will be reset. If necessary, please make a backup in advance.***

Additional software packages required for GPU devices.
- nvidia-docker2

## Instructions for use

***Please make sure to switch to the root account before starting.***

1. Compile binary executable file.

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/DistriAI main.go
```

2. Write configuration file.

Create config.yml file in the same directory as the executable file.
```
cd bin
cat << EOF > config.yml
base:
  # Mnemonics used to complete transactions
  privateKey: <your_private_key>
  # The level of privacy protection provided
  securityLevel: 0
console:
  # Directory provided for training models. default: /data/distri
  workDirectory:
  # If using a proxy, please specify the public IP address. If not using a proxy, you can ignore this.
  outerNetIP:
  # If not filled, the value of the 'nginxPost' field will be used to fill this.
  outerNetPort:
  # The port on which the nginx listens. default: 80
  nginxPost:
  # The port on which the console listens. default: 8080
  consolePost:
  # The port on which the Local server listens. default: 8088
  serverPost:
EOF
```

3. Run executable file.

```
./DistriAI node start
```