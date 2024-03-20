package cmd

import (

	// "github.com/jdgcs/ed25519/extra25519"

	// "github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/client"
	// "github.com/docker/go-connections/nat"

	"DistriAI-Node/utils"

	"github.com/urfave/cli"
	// "golang.org/x/crypto/nacl/box"
)

var DebugCommand = cli.Command{
	Name:  "debug",
	Usage: "debug.",
	Action: func(c *cli.Context) error {

		/* Debug : Force Complete */
		// machineUUID, err := machine_uuid.GetInfoMachineUUID()
		// if err != nil {
		// 	return err
		// }

		// key := config.GlobalConfig.Base.PrivateKey

		// newConfig := config.NewConfig(
		// 	key,
		// 	pattern.RPC,
		// 	pattern.WsRPC)

		// var chainInfo *chain.InfoChain
		// chainInfo, err = chain.GetChainInfo(newConfig, machineUUID)
		// if err != nil {
		// 	return err
		// }

		// var orderPlacedMetadata pattern.OrderPlacedMetadata

		// metadata := "{\"formData\":{\"taskName\":\"Computing Task - 2\",\"duration\":1},\"MachineInfo\":{\"UUID\":\"0565268338504c89ba51231d75ab4735\",\"Provider\":\"AxBoDKGYKBa54qkDusWWYgf8QXufvBKTJTQBaKyEiEzF\",\"Region\":\"China\",\"GPU\":\"1xNVIDIA GeForce GTX 1080 Ti\",\"CPU\":\"11th Gen Intel(R) Core(TM) i5-11400 @ 2.60GHz\",\"TFLOPS\":0,\"RAM\":\"15GB\",\"AvailDiskStorage\":100,\"Reliability\":\"100%\",\"CPS\":\"79.27\",\"Speed\":{\"Upload\":\"74.82 Mbit/s\",\"Download\":\"50.16 Mbit/s\"}},\"machinePublicKey\":\"9ycKDBBa3b1AZWjT6yrePdaPCypFsSSyQtsMuezDb6vv\"}"

		// err = json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
		// if err != nil {
		// 	return err
		// }

		// orderPlacedMetadata.MachineAccounts = chainInfo.ProgramDistriMachine.String()

		// chainInfo.ProgramDistriOrder = solana.MustPublicKeyFromBase58("5EAjFP6DSX9zAmwNYbmpduYho8Y4sxiZmT5Qhmrz6uVT")

		// buyer := solana.MustPublicKeyFromBase58("4F1fmZAmZ7bwQF3vz3Dv7VFJDyrkyjDyftsKHq9bTb1p")

		// distriWrapper := distri.NewDistriWrapper(chainInfo)
		// _, err = distriWrapper.OrderFailed(buyer, orderPlacedMetadata)
		// if err != nil {
		// 	logs.Error(err.Error())
		// }
		// return nil

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

		// containerID, err := docker.RunWorkspaceContainer(true, mlToken)
		// if err != nil {
		// 	return err
		// }
		// logs.Normal(fmt.Sprintf("containerID: %v", containerID))

		// server.StartServer(config.GlobalConfig.Console.ServerPost)

		urls := []utils.DownloadURL{
			{
				URL:      "https://distriai.s3.ap-northeast-2.amazonaws.com/model/Bv3qEmRjPn3z7bB3JynCoXJmopcNM8PGa6ASxPCi7bY/animagine-xl-3.0/pytorch_model.bin",
				Checksum: "",
			},
			{
				URL:      "https://distriai.s3.ap-northeast-2.amazonaws.com/model/Bv3qEmRjPn3z7bB3JynCoXJmopcNM8PGa6ASxPCi7bY/animagine-xl-3.0/pytorch_model1.bin",
				Checksum: "",
			},
			{
				URL:      "https://distriai.s3.ap-northeast-2.amazonaws.com/model/Bv3qEmRjPn3z7bB3JynCoXJmopcNM8PGa6ASxPCi7bY/animagine-xl-3.0/pytorch_model2.bin",
				Checksum: "",
			},
		}

		// name, err := utils.GetFilenameFromURL(urls[0].URL)
		// if err != nil {
		// 	return err
		// }
		// logs.Normal(name)
		err := utils.DownloadFiles("/data/download", urls)
		if err != nil {
			return err
		}
		return nil
	},
}
