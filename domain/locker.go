package domain

type LockRequest struct {
	Key      string `validate:"required"`
	TTLInSec int    `validate:"required"`
}

type UnLockRequest struct {
	Key     string `validate:"required"`
	LockKey string `validate:"required"`
}

type LockResponse struct {
	// nolint: tagliatelle
	LockKey string `json:",omitempty"`
}
