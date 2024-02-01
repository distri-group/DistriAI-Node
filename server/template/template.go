package template

const (
	WORKSPACE = "/workspace"
	TOKEN     = "/token"
)

const (
	ORDER = "/order"
	RT    = "/rt"
)

type RequestToken struct {
	Signature string `json:"signature"`
}
