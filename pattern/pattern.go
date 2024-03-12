package pattern

import "github.com/gagliardetto/solana-go/rpc"

// RPC is the url of the node
const RPC = rpc.DevNet_RPC
const WsRPC = rpc.DevNet_WS

const PROGRAM_DISTRI_ID = "8WxFh7ZtvTP1587yvyweoRMt41hCHMsmM48PTbZKM4tY"

const DIST_TOKEN_ID = "896KfVVY6VRGQs1d9CKLnKUEgXXCCJcEEg7LwSK84vWE"

const NO_GPU = "No GPU"

// docer
const (
	DOCKER_GROUP = "distrigroup"
)

// docker: score image
const (
	SCORE_IMAGE     = "ml-device-score"
	SCORE_TAGS      = "v0.0.1"
	SCORE_CONTAINER = "ml-device-score"
	SCORE_NAME      = DOCKER_GROUP + "/" + SCORE_IMAGE + ":" + SCORE_TAGS
)

// docker: ml-workspace image
const (
	ML_WORKSPACE_IMAGE     = "ml-workspace-gpu"
	ML_WORKSPACE_TAGS      = "0.3.4"
	ML_WORKSPACE_CONTAINER = "ml-workspace-gpu"
	ML_WORKSPACE_NAME      = DOCKER_GROUP + "/" + ML_WORKSPACE_IMAGE + ":" + ML_WORKSPACE_TAGS
)

// docker: ml-workspace-gpu image
const (
	ML_WORKSPACE_GPU_IMAGE     = "ml-workspace-gpu"
	ML_WORKSPACE_GPU_TAGS      = "0.3.4"
	ML_WORKSPACE_GPU_CONTAINER = "ml-workspace-gpu"
	ML_WORKSPACE_GPU_NAME      = DOCKER_GROUP + "/" + ML_WORKSPACE_GPU_IMAGE + ":" + ML_WORKSPACE_TAGS
)

// DOT is "." character
const DOT = "."

const (
	// HASHRATE_MARKET is a module about DeOSS
	HASHRATE_MARKET = "HashrateMarket"
)

// Extrinsic
const (
	// TX_HASHRATE_MARKET_REGISTER
	TX_HASHRATE_MARKET_ORDER_START = HASHRATE_MARKET + DOT + "order_start"

	TX_HASHRATE_MARKET_REGISTER = HASHRATE_MARKET + DOT + "add_machine"

	TX_HASHRATE_MARKET_ORDER_COMPLETED = HASHRATE_MARKET + DOT + "order_completed"

	TX_HASHRATE_MARKET_ORDER_FAILED = HASHRATE_MARKET + DOT + "order_failed"

	TX_HASHRATE_MARKET_REMOVE_MACHINE = HASHRATE_MARKET + DOT + "remove_machine"

	TX_HASHRATE_MARKET_SUBMIT_TASK = HASHRATE_MARKET + DOT + "submit_task"
)

type MachineUUID [16]byte
type TaskUUID [16]byte

type OrderPlacedMetadata struct {
	FormData        FormData    `json:"formData"`
	MachineInfo     MachineInfo `json:"MachineInfo"`
	MachineAccounts string      `json:"MachineAccounts"`
}

type MachineInfo struct {
	UUID             string    `json:"UUID"`
	Provider         string    `json:"Provider"`
	Region           string    `json:"Region"`
	GPU              string    `json:"GPU"`
	CPU              string    `json:"CPU"`
	TFLOPS           float32   `json:"TFLOPS"`
	RAM              string    `json:"RAM"`
	AvailDiskStorage int32     `json:"AvailDiskStorage"`
	Reliability      string    `json:"Reliability"`
	CPS              string    `json:"CPS"`
	Speed            SpeedInfo `json:"Speed"`
}

type SpeedInfo struct {
	Upload   string `json:"Upload"`
	Download string `json:"Download"`
}

type FormData struct {
	TaskName  string `json:"taskName"`
	Duration  int    `json:"duration"`
	BuyTime   string `json:"buyTime"`
	OrderTime string `json:"orderTime"`
}

type TaskMetadata struct {
}
