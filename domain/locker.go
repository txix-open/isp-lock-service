package domain

type LockRequest struct {
	Key      string `valid:"required"`
	TTLInSec int    `valid:"required"`
}

type UnLockRequest struct {
	Key     string `valid:"required"`
	LockKey string `valid:"required"`
}

type LockResponse struct {
	LockKey string `json:",omitempty"`
}
