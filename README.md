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

Additional software packages required for GPU devices.
- nvidia-container-toolkit
- nvidia-docker2

## Instructions for use
1. Compile binary executable file.

```
go build -o bin/DistriAI main.go
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
  # If not filled, the value of the 'innerNetPort' field will be used to fill this.
  outerNetPort:
  # The port on which the console listens. default: 8080
  innerNetPort:
EOF
```

3. Run executable file.

```
./DistriAI node start
```