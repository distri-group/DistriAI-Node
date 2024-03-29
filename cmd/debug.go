package cmd

import (

	// "github.com/jdgcs/ed25519/extra25519"

	// "github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/client"
	// "github.com/docker/go-connections/nat"

	"DistriAI-Node/chain"
	"DistriAI-Node/chain/distri"
	"DistriAI-Node/config"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/pattern"
	logs "DistriAI-Node/utils/log_utils"
	"encoding/json"

	"github.com/gagliardetto/solana-go"
	"github.com/urfave/cli"
	// "golang.org/x/crypto/nacl/box"
)

var DebugCommand = cli.Command{
	Name:  "debug",
	Usage: "debug.",
	Action: func(c *cli.Context) error {

		/* Debug : Force Complete */
		machineUUID, err := machine_uuid.GetInfoMachineUUID()
		if err != nil {
			return err
		}

		key := config.GlobalConfig.Base.PrivateKey

		newConfig := config.NewConfig(
			key,
			pattern.RPC,
			pattern.WsRPC)

		var chainInfo *chain.InfoChain
		chainInfo, err = chain.GetChainInfo(newConfig, machineUUID)
		if err != nil {
			return err
		}

		var orderPlacedMetadata pattern.OrderPlacedMetadata

		metadata := "{\"formData\":{\"duration\":1,\"taskName\":\"Computing Task-33\"},\"MachineInfo\":{\"UUID\":\"ec29d2a3323b6812b091d70610d2e17b\",\"Provider\":\"AxBoDKGYKBa54qkDusWWYgf8QXufvBKTJTQBaKyEiEzF\",\"Region\":\"United States\",\"GPU\":\"1xNVIDIA A10G\",\"CPU\":\"AMD EPYC 7R32\",\"Tflops\":31.52,\"RAM\":\"15GB\",\"AvailDiskStorage\":100,\"Reliability\":\"--\",\"CPS\":\"53.72\",\"Speed\":{\"Upload\":\"1258.21 Mbit/s\",\"Download\":\"1543.24 Mbit/s\"},\"MaxDuration\":10,\"Price\":1,\"IP\":\"3.235.182.86\",\"Port\":\"80\"},\"OrderInfo\":{\"Model\":1,\"Intent\":\"deploy\",\"DownloadURL\":[\"https://distriai.s3.ap-northeast-2.amazonaws.com/model/Bv3qEmRjPn3z7bB3JynCoXJmopcNM8PGa6ASxPCi7bY/animagine-xl-3.0/debug-deploy.py\"]},\"machinePublicKey\":\"AWufzFBiJTHrBVbZwTLR2Aedg6v6C3bqJNy1LaAwZnRe\"}"

		err = json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
		if err != nil {
			return err
		}

		orderPlacedMetadata.MachineAccounts = chainInfo.ProgramDistriMachine.String()

		chainInfo.ProgramDistriOrder = solana.MustPublicKeyFromBase58("9D1oF3rYzZpJnQTjzeMyC7Rowiujd6PzMpD74JLH7dY6")

		buyer := solana.MustPublicKeyFromBase58("AxBoDKGYKBa54qkDusWWYgf8QXufvBKTJTQBaKyEiEzF")

		distriWrapper := distri.NewDistriWrapper(chainInfo)
		_, err = distriWrapper.OrderFailed(buyer, orderPlacedMetadata)
		if err != nil {
			logs.Error(err.Error())
		}
		return nil

		// db, err := dbutils.NewDB()
		// if err != nil {
		// 	logs.Error(err.Error())
		// }
		// db.Delete([]byte("buyer"))
		// db.Delete([]byte("token"))
		// db.Close()

		// nginx.StopNginx()

		/* Debug : solana Sign */
		// private := solana.MustPrivateKeyFromBase58("2hpaRRjxuzhiHJ6ggimgGHi8jCMgG3MRXTUAGm7XNrztXpzKvJJ4fqAfgzb4YNnT9LDvYMfh4GzFh5NBojUcCXaB")
		// public := private.PublicKey()
		// logs.Normal(fmt.Sprintf("publicKey: %v", public))
		// msg := "workspace/token/" + public.String()
		// signature, err := private.Sign([]byte(msg))
		// if err != nil {
		// 	return err
		// }
		// logs.Normal(fmt.Sprintf("signature: %v", signature.String()))
		// if public.Verify([]byte(msg), signature) {
		// 	logs.Normal("Verify success")
		// } else {
		// 	logs.Error("Verify failed")
		// }

		/* Debug : docker space */
		// sizeLimitGB := 10

		// result, err := utils.CompareSpaceWithDocker(sizeLimitGB)
		// if err != nil {
		// 	logs.Error(fmt.Sprintf("比较空间大小时发生错误: %v\n", err))
		// 	return nil
		// }

		// if result {
		// 	logs.Normal("Docker 镜像目录的可用空间充足")
		// } else {
		// 	logs.Normal("Docker 镜像目录的可用空间不足")
		// }
		// return nil

		/* Debug : nginx ssl */
		// if err := nginx.StartNginx(
		// 	config.GlobalConfig.Console.NginxPort,
		// 	config.GlobalConfig.Console.ConsolePost,
		// 	config.GlobalConfig.Console.ServerPost); err != nil {
		// 	logs.Error(err.Error())
		// 	return nil
		// }

		// mlToken, err := utils.GenerateRandomString(16)
		// if err != nil {
		// 	logs.Error(err.Error())
		// 	return nil
		// }

		// db, err := dbutils.NewDB()
		// if err != nil {
		// 	logs.Error(err.Error())
		// 	return nil
		// }
		// db.Update([]byte("buyer"), []byte(public.String()))
		// db.Update([]byte("token"), []byte(mlToken))
		// db.Close()
		// logs.Normal(fmt.Sprintf("From buyer: %v ; mlToken: %v", public, mlToken))

		// containerID, err := docker.RunWorkspaceContainer(true, "mlToken")
		// if err != nil {
		// 	return err
		// }
		// logs.Normal(fmt.Sprintf("containerID: %v", containerID))

		// server.StartServer(config.GlobalConfig.Console.ServerPost)

		// urls := []string{
		// 	"https://distriai.s3.ap-northeast-2.amazonaws.com/model/Bv3qEmRjPn3z7bB3JynCoXJmopcNM8PGa6ASxPCi7bY/animagine-xl-3.0/debug-deploy.py",
		// }

		// id, err := docker.RunDeployContainer(true, urls)
		// if err != nil {
		// 	return err
		// }
		
		// logs.Normal(fmt.Sprintf("containerID: %v", id))
	},
}