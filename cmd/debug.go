package cmd

import (
	"DistriAI-Node/chain"
	"DistriAI-Node/chain/distri"
	"DistriAI-Node/config"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/pattern"
	"encoding/json"

	// "github.com/jdgcs/ed25519/extra25519"
	"github.com/gagliardetto/solana-go"
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

		metadata := "{\"formData\":{\"taskName\":\"Computing Task - 7\",\"duration\":1},\"machinePublicKey\":\"9fMCEWm5Y6tqn9UseHrfPmokWAbjJBooK5mjPanAYgbb\"}"

		err = json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
		if err != nil {
			return err
		}

		orderPlacedMetadata.MachineAccounts = chainInfo.ProgramDistriMachine.String()

		chainInfo.ProgramDistriOrder = solana.MustPublicKeyFromBase58("9fMCEWm5Y6tqn9UseHrfPmokWAbjJBooK5mjPanAYgbb")

		distriWrapper := distri.NewDistriWrapper(chainInfo)
		_, err = distriWrapper.OrderCompleted(orderPlacedMetadata, false)
		if err != nil {
			return err
		}
		return nil
	},
}
