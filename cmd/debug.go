package cmd

import (

	// "github.com/jdgcs/ed25519/extra25519"

	// "github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/client"
	// "github.com/docker/go-connections/nat"

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

		// metadata := "{\"formData\":{\"duration\":2,\"taskName\":\"Computing Task-0\"},\"MachineInfo\":{\"UUID\":\"ec29d2a3323b6812b091d70610d2e17b\",\"Provider\":\"AxBoDKGYKBa54qkDusWWYgf8QXufvBKTJTQBaKyEiEzF\",\"Region\":\"United States\",\"GPU\":\"1xNVIDIA A10G\",\"CPU\":\"AMD EPYC 7R32\",\"Tflops\":31.52,\"RAM\":\"15GB\",\"AvailDiskStorage\":100,\"Reliability\":\"--\",\"CPS\":\"53.05\",\"Speed\":{\"Upload\":\"1295.07 Mbit/s\",\"Download\":\"1542.62 Mbit/s\"},\"MaxDuration\":10,\"Price\":1,\"IP\":\"3.235.182.86\",\"Port\":\"80\"},\"OrderInfo\":{\"Model\":\"\",\"Intent\":\"\",\"DownloadURL\":[]},\"machinePublicKey\":\"33fYCQYz2HteuUyVNh51zAjfRNyVmT2t5JGGbYNFS6nc\"}"

		// err = json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
		// if err != nil {
		// 	return err
		// }

		// orderPlacedMetadata.MachineAccounts = chainInfo.ProgramDistriMachine.String()

		// chainInfo.ProgramDistriOrder = solana.MustPublicKeyFromBase58("9K1Mdr7d455mGZrAvdBUArjrpS4ujhCZB72mNqeFdyWV")

		// buyer := solana.MustPublicKeyFromBase58("4F1fmZAmZ7bwQF3vz3Dv7VFJDyrkyjDyftsKHq9bTb1p")

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
		// 	logs.Error(fmt.Sprintf("An error occurred while comparing disk space sizes: %v\n", err))
		// 	return nil
		// }

		// if result {
		// 	logs.Normal("There is sufficient free space in the Docker directory")
		// } else {
		// 	logs.Normal("Insufficient free space in the Docker directory")
		// }
		// return nil

		/* Debug : GetTransaction */
		// endpoint := rpc.DevNet_RPC
		// client := rpc.New(endpoint)

		// txSig := solana.MustSignatureFromBase58("4WYFFjWh8YhPwsd1gBK6zo4TNzDEN34ufY4eQFUrqn61uJ5Z7EZpjFnLwvZYLdqaeSpdHwA91fdzM8uSVZgfKuUG")
		// {
		//   out, err := client.GetTransaction(
		// 	context.TODO(),
		// 	txSig,
		// 	&rpc.GetTransactionOpts{
		// 	  Encoding: solana.EncodingBase64,
		// 	},
		//   )
		//   if err != nil {
		// 	panic(err)
		//   }
		//   logs.Normal("111111111111111")
		//   spew.Dump(out)
		//   logs.Normal("222222222222222")
		//   spew.Dump(out.Transaction.GetBinary())

		//   decodedTx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
		//   if err != nil {
		// 	panic(err)
		//   }
		//   logs.Normal("333333333333333")
		//   spew.Dump(decodedTx)
		// }
		// {
		//   out, err := client.GetTransaction(
		// 	context.TODO(),
		// 	txSig,
		// 	nil,
		//   )
		//   if err != nil {
		// 	panic(err)
		//   }
		//   logs.Normal("555555555555555")
		//   spew.Dump(out)
		//   logs.Normal("666666666666666")
		//   spew.Dump(out.Transaction.GetTransaction())
		// }
		return nil
	},
}
