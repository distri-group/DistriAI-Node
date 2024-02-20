package cmd

import (

	// "github.com/jdgcs/ed25519/extra25519"

	// "github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/client"
	// "github.com/docker/go-connections/nat"

	docker_utils "DistriAI-Node/docker/utils"
	logs "DistriAI-Node/utils/log_utils"
	"context"

	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	// "golang.org/x/crypto/nacl/box"
)

var DebugCommand = cli.Command{
	Name:  "debug",
	Usage: "debug.",
	Action: func(c *cli.Context) error {

		/* Debug : Asymmetric Encryption */
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

		// logs.Normal(fmt.Sprintf("publicKey : %v", chainInfo.Wallet.Wallet.PublicKey()))
		// logs.Normal(fmt.Sprintf("privateKey : %v", (*[64]byte)(chainInfo.Wallet.Wallet.PrivateKey)))

		// var curve25519PrivateKey [32]byte
		// extra25519.PrivateKeyToCurve25519(&curve25519PrivateKey, (*[64]byte)(chainInfo.Wallet.Wallet.PrivateKey))

		// var curve25519PublicKey [32]byte
		// if !extra25519.PublicKeyToCurve25519(&curve25519PublicKey, (*[32]byte)(chainInfo.Wallet.Wallet.PublicKey().Bytes())) {
		// 	logs.Error("error in public key conversion")
		// }

		// logs.Normal(fmt.Sprintf("curve25519PrivateKey: %v", curve25519PrivateKey))
		// logs.Normal(fmt.Sprintf("curve25519PublicKey: %v", curve25519PublicKey))

		// message := []byte("Hello, World! GGH")
		// var nonce [24]byte
		// _, _ = rand.Read(nonce[:])
		// fmt.Println(nonce)
		// encrypted := box.Seal(nonce[:], message, &nonce, &curve25519PublicKey, &curve25519PrivateKey)
		// fmt.Println(encrypted)

		// nonce = [24]byte{185, 216, 206, 50, 246, 77, 250,
		// 	134, 168, 64, 29, 254, 218, 43,
		// 	168, 15, 66, 214, 250, 9, 166,
		// 	47, 27, 34}
		// encrypted = []byte{185, 216, 206, 50, 246, 77, 250,
		// 	134, 168, 64, 29, 254, 218, 43,
		// 	168, 15, 66, 214, 250, 9, 166,
		// 	47, 27, 34, 207, 27, 1, 146, 61, 136, 32, 95,
		// 	213, 194, 134, 133, 139, 101, 32, 101,
		// 	51, 241, 127, 56, 73, 71, 98, 146,
		// 	30, 55, 90, 89, 114, 98, 49}

		// decrypted, ok := box.Open(nil, encrypted[24:], &nonce, &curve25519PublicKey, &curve25519PrivateKey)
		// if !ok {
		// 	fmt.Println("decryption error")
		// 	return nil
		// }

		// fmt.Println(string(decrypted)) // 输出: Hello, World!

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

		// metadata := "{\"formData\":{\"taskName\":\"Computing Task - 32\",\"duration\":1},\"MachineInfo\":{\"UUID\":\"0565268338504c89ba51231d75ab4735\",\"Provider\":\"AxBoDKGYKBa54qkDusWWYgf8QXufvBKTJTQBaKyEiEzF\",\"Region\":\"China\",\"GPU\":\"1xNVIDIA GeForce GTX 1080 Ti\",\"CPU\":\"11th Gen Intel(R) Core(TM) i5-11400 @ 2.60GHz\",\"TFLOPS\":0,\"RAM\":\"15GB\",\"AvailDiskStorage\":100,\"Reliability\":\"--\",\"CPS\":\"78.97\",\"Speed\":{\"Upload\":\"76.00 Mbit/s\",\"Download\":\"23.64 Mbit/s\"}},\"machinePublicKey\":\"Go7DjYCFcKXZ1AUdWW3wq9yqQCMBJDL4Vsvu2qEdcyAv\"}"

		// err = json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
		// if err != nil {
		// 	return err
		// }

		// orderPlacedMetadata.MachineAccounts = chainInfo.ProgramDistriMachine.String()

		// chainInfo.ProgramDistriOrder = solana.MustPublicKeyFromBase58("AymhSXBLEAF61LXNnN874kDEP66zR19WVg6cf4LfYfbb")

		// buyer := solana.MustPublicKeyFromBase58("ExCX1FnGPjYAbXREqACWp7wSWe2jFXon6pJXTKTxsn4k")

		// distriWrapper := distri.NewDistriWrapper(chainInfo)
		// _, err = distriWrapper.OrderFailed(buyer, orderPlacedMetadata)
		// if err != nil {
		// 	logs.Error(err.Error())
		// }

		// db, err := dbutils.NewDB()
		// if err != nil {
		// 	logs.Error(err.Error())
		// }
		// db.Delete([]byte("buyer"))
		// db.Delete([]byte("token"))
		// db.Close()

		// nginx.StopNginx()

		/* Dedug : ml-workspace */
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
		// db.Update([]byte("buyer"), []byte("ExCX1FnGPjYAbXREqACWp7wSWe2jFXon6pJXTKTxsn4k"))
		// db.Update([]byte("token"), []byte(mlToken))
		// db.Close()

		// config.GlobalConfig.Console.WorkDirectory = "/data/debug"
		// _, err = docker.RunWorkspaceContainer(true, mlToken)
		// if err != nil {
		// 	logs.Error(fmt.Sprintln("RunWorkspaceContainer error: ", err))
		// 	return nil
		// }

		/* Dedug : nginx */
		// err = nginx.StartNginx(config.GlobalConfig.Console.NginxPort,
		// 	config.GlobalConfig.Console.ConsolePost,
		// 	config.GlobalConfig.Console.ServerPost)
		// if err != nil {
		// 	logs.Error(fmt.Sprintln("StartNginx error: ", err))
		// 	return nil
		// }

		// if err := server.StartServer(config.GlobalConfig.Console.ServerPost); err != nil {
		// 	logs.Error(err.Error())
		// 	return nil
		// }

		/* Dedug : solana Sign */
		// private := solana.MustPrivateKeyFromBase58("3wXVb6mVr5UHsetWYZQkGW3er86M3hXXgjw4LFcxAugRg9EQcZHVcW2hWeuTmXahgyQgdRZfuY2XRPaKLm4v2ywz")
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

		/* Debug : docker container */
		ctx, cancel := context.WithCancel(context.Background())
		cli, _ := client.NewClientWithOpts(client.FromEnv)
		defer cancel()

		if is, id := docker_utils.ContainerExists(ctx, cli, "debug-workspace-ggh"); is {
			logs.Normal("Container exists")
			logs.Normal(id)
		} else {
			logs.Normal("Container does not exist")
		}
		return nil
	},
}
