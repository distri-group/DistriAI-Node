package chain

import (
	"DistriAI-Node/chain/conn"
	"DistriAI-Node/chain/wallet"
	"DistriAI-Node/config"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/pattern"
	logs "DistriAI-Node/utils/log_utils"
	"encoding/hex"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type InfoChain struct {
	Conn                 *conn.Conn
	Wallet               *wallet.Wallet
	ProgramDistriID      solana.PublicKey
	ProgramDistriMachine solana.PublicKey
	ProgramDistriOrder   solana.PublicKey
}

// GetChainInfo 创建conn\wallet信息并返回ChainInfo
func GetChainInfo(cfg *config.SolanaConfig, machineUUID machine_uuid.MachineUUID) (*InfoChain, error) {
	newConn, err := conn.NewConn(cfg)
	if err != nil {
		return nil, err
	}

	wallet, err := wallet.InitWallet(cfg)
	if err != nil {
		return nil, err
	}

	programID := solana.MustPublicKeyFromBase58(pattern.PROGRAM_DISTRI_ID)

	var byteUUID pattern.MachineUUID
	b, err := hex.DecodeString(string(machineUUID))
	if err != nil {
		panic(err)
	}
	copy(byteUUID[:], b[:16])

	seedMachine := [][]byte{
		[]byte(pattern.DISTRI_SEED),
		wallet.Wallet.PublicKey().Bytes(),
		[]byte(byteUUID[:]),
	}

	machineAccount, _, err := solana.FindProgramAddress(
		seedMachine,
		programID,
	)
	if err != nil {
		return nil, err
	}
	logs.Normal(fmt.Sprintf("machineAccount : %v", machineAccount.String()))

	chainInfo := &InfoChain{
		Conn:                 newConn,
		Wallet:               wallet,
		ProgramDistriID:      programID,
		ProgramDistriMachine: machineAccount,
	}

	return chainInfo, nil
}
