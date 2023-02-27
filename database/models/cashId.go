package models

var StatusCodes = map[string]int{
	"SUCCESSFUL": 0,

	"REQUEST_BROKEN":           100,
	"REQUEST_MISSING_SCHEME":   111,
	"REQUEST_MISSING_DOMAIN":   112,
	"REQUEST_MISSING_NONCE":    113,
	"REQUEST_MALFORMED_SCHEME": 121,
	"REQUEST_MALFORMED_DOMAIN": 122,
	"REQUEST_INVALID_DOMAIN":   131,
	"REQUEST_INVALID_NONCE":    132,
	"REQUEST_ALTERED":          141,
	"REQUEST_EXPIRED":          142,
	"REQUEST_CONSUMED":         143,

	"RESPONSE_BROKEN":              200,
	"RESPONSE_MISSING_REQUEST":     211,
	"RESPONSE_MISSING_ADDRESS":     212,
	"RESPONSE_MISSING_SIGNATURE":   213,
	"RESPONSE_MISSING_METADATA":    214,
	"RESPONSE_MALFORMED_ADDRESS":   221,
	"RESPONSE_MALFORMED_SIGNATURE": 222,
	"RESPONSE_MALFORMED_METADATA":  223,
	"RESPONSE_INVALID_METHOD":      231,
	"RESPONSE_INVALID_ADDRESS":     232,
	"RESPONSE_INVALID_SIGNATURE":   233,
	"RESPONSE_INVALID_METADATA":    234,

	"SERVICE_BROKEN":                 300,
	"SERVICE_ADDRESS_DENIED":         311,
	"SERVICE_ADDRESS_REVOKED":        312,
	"SERVICE_ACTION_DENIED":          321,
	"SERVICE_ACTION_UNAVAILABLE":     322,
	"SERVICE_ACTION_NOT_IMPLEMENTED": 323,
	"SERVICE_INTERNAL_ERROR":         331,
}

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
