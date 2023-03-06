package models

var ServiceActions = []string{"auth", "login", "sign", "register", "ticket", "claimtx?", "claimaddr?"}
var UserActions = []string{"delete", "logout", "revoke", "update"}

type CashIDRequest struct {
	Intent   string
	Domain   string
	Path     string
	Action   string
	Data     map[string]string
	Required map[string]interface{}
	Optional map[string]interface{}
	Nonce    string
}

type StatusConfirmation struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CashId struct {
	StatusConfirmation `json:"status-confirmation"`
}
type Payload struct {
	Address   string                 `json:"address" binding:"required"`
	Signature string                 `json:"signature" binding:"required"`
	Request   string                 `json:"request" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata"`
}
