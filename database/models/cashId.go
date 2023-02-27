package models

type StatusConfirmation struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CashId struct {
	StatusConfirmation `json:"status-confirmation"`
}
