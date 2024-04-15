package chain

import (
	"DistriAI-Node/chain/conn"
	"DistriAI-Node/chain/wallet"
	"DistriAI-Node/config"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/pattern"
	"DistriAI-Node/utils"
	logs "DistriAI-Node/utils/log_utils"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// InfoChain encapsulates basic information about the chain
type InfoChain struct {
	Conn                 *conn.Conn
	Wallet               *wallet.Wallet
	ProgramDistriID      solana.PublicKey
	ProgramDistriMachine solana.PublicKey
	ProgramDistriOrder   solana.PublicKey
}

// GetChainInfo returns *Infochain and error when the connection fails
func GetChainInfo(cfg *config.SolanaConfig, machineUUID machine_uuid.MachineUUID) (*InfoChain, error) {
	newConn, err := conn.NewConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("> conn.NewConn: %v", err)
	}

	wallet, err := wallet.InitWallet(cfg)
	if err != nil {
		return nil, fmt.Errorf("> wallet.InitWallet: %v", err)
	}

	programID := solana.MustPublicKeyFromBase58(pattern.PROGRAM_DISTRI_ID)

	seedMachine := utils.GenMachine(wallet.Wallet.PublicKey(), machineUUID)

	machineAccount, _, err := solana.FindProgramAddress(
		seedMachine,
		programID,
	)
	if err != nil {
		return nil, fmt.Errorf("> FindProgramAddress: %v", err)
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
