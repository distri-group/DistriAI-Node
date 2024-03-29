package subscribe

import (
	"DistriAI-Node/chain"
	"DistriAI-Node/chain/distri/distri_ai"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/utils"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type WrapperSubscribe struct {
	*chain.InfoChain
}

func (chain *WrapperSubscribe) SubscribeProgram(MachineUUID machine_uuid.MachineUUID) (distri_ai.Order, error) {

	var order distri_ai.Order

	sub, err := chain.Conn.WsClient.ProgramSubscribeWithOpts(
		chain.ProgramDistriID,
		rpc.CommitmentFinalized,
		solana.EncodingBase64Zstd,
		nil,
	)
	if err != nil {
		return order, fmt.Errorf("error subscribing to program: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return order, err
		}
		borshDec := bin.NewBorshDecoder(got.Value.Account.Data.GetBinary())

		err = order.UnmarshalWithDecoder(borshDec)
		if err != nil {
			continue
		} else {
			uuid, err := utils.ParseMachineUUID(string(MachineUUID))
			if err != nil {
				return order, fmt.Errorf("error parsing uuid: %v", err)
			}
			if order.Seller.Equals(chain.Wallet.Wallet.PublicKey()) && order.MachineId == uuid {
				chain.ProgramDistriOrder = got.Value.Pubkey
				return order, nil
			}
		}
	}
}

func (chain *WrapperSubscribe) SubscribeAccount() (distri_ai.Order, error) {

	var order distri_ai.Order

	sub, err := chain.Conn.WsClient.AccountSubscribeWithOpts(
		chain.ProgramDistriID,
		rpc.CommitmentFinalized,
		solana.EncodingBase64Zstd,
	)
	if err != nil {
		return order, fmt.Errorf("> AccountSubscribeWithOpts: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return order, fmt.Errorf("> Recv: %v", err)
		}
		borshDec := bin.NewBorshDecoder(got.Value.Account.Data.GetBinary())

		err = order.UnmarshalWithDecoder(borshDec)
		if err != nil {
			continue
		}
		return order, nil
	}
}

func NewSubscribeWrapper(info *chain.InfoChain) *WrapperSubscribe {
	return &WrapperSubscribe{info}
}
