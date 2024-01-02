package conn

import (
	"DistriAI-Node/config"
	"context"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

type Conn struct {
	RpcClient *rpc.Client
	WsClient  *ws.Client
}

func NewConn(cfg *config.SolanaConfig) (*Conn, error) {

	rpcClient := rpc.New(cfg.RPC)
	wsClient, err := ws.Connect(context.Background(), cfg.WsRPC)
	if err != nil {
		return nil, err
	}

	conn := &Conn{
		RpcClient: rpcClient,
		WsClient:  wsClient,
	}

	return conn, nil
}
