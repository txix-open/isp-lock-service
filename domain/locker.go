package domain

import "time"

type LockRequest struct {
	Key      string        `valid:"required"`
	TTLInSec time.Duration `valid:"required"`
}

type UnLockRequest struct {
	Key     string `valid:"required"`
	LockKey string `valid:"required"`
}

type LockResponse struct {
	LockKey string `json:",omitempty"`
}
