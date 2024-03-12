package conn

import (
	"DistriAI-Node/config"
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

//  Connection to Solana Blockchain Nodes
type Conn struct {
	RpcClient *rpc.Client
	WsClient  *ws.Client
}
// Receive a SolanaConfig configuration object and return an initialized Conn object.
func NewConn(cfg *config.SolanaConfig) (*Conn, error) {

	rpcClient := rpc.New(cfg.RPC)
	wsClient, err := ws.Connect(context.Background(), cfg.WsRPC)
	if err != nil {
		return nil, fmt.Errorf("> ws.Connect: %v", err)
	}

	conn := &Conn{
		RpcClient: rpcClient,
		WsClient:  wsClient,
	}

	return conn, nil
}
