package pattern

import "github.com/gagliardetto/solana-go/rpc"

// RPC is the url of the node
const RPC = rpc.DevNet_RPC  // 开发链
const WsRPC = rpc.DevNet_WS // 开发链

// DOT is "." character
const DOT = "."

const (
	// HASHRATE_MARKET is a module about DeOSS
	HASHRATE_MARKET = "HashrateMarket"
)

// Extrinsic
const (

	// TX_HASHRATE_MARKET_REGISTER
	TX_HASHRATE_MARKET_REGISTER = HASHRATE_MARKET + DOT + "add_machine"

	TX_HASHRATE_MARKET_ORDER_COMPLETED = HASHRATE_MARKET + DOT + "order_completed"

	TX_HASHRATE_MARKET_ORDER_FAILED = HASHRATE_MARKET + DOT + "order_failed"

	TX_HASHRATE_MARKET_REMOVE_MACHINE = HASHRATE_MARKET + DOT + "remove_machine"
)

type MachineUUID [16]byte

type OrderPlacedMetadata struct {
	// MachineInfo MachineInfo `json:"machineInfo"`
	FormData        FormData `json:"formData"`
	MachineAccounts string   `json:"MachineAccounts"` // 机器存储地址
}

type MachineInfo struct {
	Id             int      `json:"Id"`
	Owner          string   `json:"Owner"`
	Uuid           string   `json:"Uuid"`
	Metadata       Metadata `json:"Metadata"`
	Status         int      `json:"Status"`
	Price          int      `json:"Price"`
	MaxDuration    int      `json:"MaxDuration"`
	Disk           int      `json:"Disk"`
	CompletedCount int      `json:"CompletedCount"`
	FailedCount    int      `json:"FailedCount"`
	Score          float32  `json:"Score"`
	Gpu            string   `json:"Gpu"`
	GpuCount       int      `json:"GpuCount"`
	Region         string   `json:"Region"`
	Tflops         float32  `json:"Tflops"`
	Addr           string   `json:"Addr"`
	UuidShort      string   `json:"UuidShort"`
	Cpu            string   `json:"Cpu"`
	RAM            string   `json:"RAM"`
	AvailHardDrive string   `json:"AvailHardDrive"`
	UploadSpeed    string   `json:"UploadSpeed"`
	DownloadSpeed  string   `json:"DownloadSpeed"`
	Reliability    string   `json:"Reliability"`
	TFLOPS         float32  `json:"TFLOPS"`
}

type Metadata struct {
	MachineUUID  string       `json:"MachineUUID"`
	Addr         string       `json:"Addr"`
	CPUInfo      CPUInfo      `json:"CPUInfo"`
	DiskInfo     DiskInfo     `json:"DiskInfo"`
	Score        float32      `json:"Score"`
	InfoMemory   InfoMemory   `json:"InfoMemory"`
	GPUInfo      GPUInfo      `json:"GPUInfo"`
	LocationInfo LocationInfo `json:"LocationInfo"`
	SpeedInfo    SpeedInfo    `json:"SpeedInfo"`
	InfoFlop     InfoFlop     `json:"InfoFlop"`
}

type CPUInfo struct {
	ModelName string  `json:"ModelName"`
	Cores     int     `json:"Cores"`
	Mhz       float32 `json:"Mhz"`
}

type DiskInfo struct {
	Path       string  `json:"Path"`
	TotalSpace float32 `json:"TotalSpace"`
}

type InfoMemory struct {
	RAM float32 `json:"RAM"`
}

type GPUInfo struct {
	Model  string `json:"Model"`
	Number int    `json:"Number"`
}

type LocationInfo struct {
	Country string `json:"Country"`
	Region  string `json:"Region"`
	City    string `json:"City"`
}

type SpeedInfo struct {
	Download string `json:"Download"`
	Upload   string `json:"Upload"`
}

type InfoFlop struct {
	Flops float32 `json:"Flops"`
}

type FormData struct {
	TaskName     string `json:"taskName"`
	ImageName    string `json:"imageName"`
	ImageTag     string `json:"imageTag"`
	Libery       string `json:"libery"`
	Model        string `json:"model"`
	DataUrl      string `json:"dataUrl"`
	Iters        string `json:"iters"`
	Batchsize    string `json:"batchsize"`
	Rate         string `json:"rate"`
	Duration     int    `json:"duration"`
	LibType      string `json:"libType"`
	BuyTime      string `json:"buyTime"`
	OrderTime    string `json:"orderTime"`
	ModelUrl     string `json:"modelUrl"`
	CompleteTime string `json:"completeTime"`
	Evaluate     string `json:"evaluate"`
}

// datasets
const (
	DATASETS_FOLDER = "./server"
	ZIP_NAME        = "/datasets.zip"
)

const PROGRAM_DISTRI_ID = "HF4aT6sho2zTySB8nEeN5ThMvDGtGVRrH3jeBvxFNxit"

const DIST_TOKEN_ID = "896KfVVY6VRGQs1d9CKLnKUEgXXCCJcEEg7LwSK84vWE"

const DISTRI_SEED = "machine"

const DISTRI_ORDER = "order"

const DISTRI_VAULT = "vault"

const NO_GPU = "No GPU"
