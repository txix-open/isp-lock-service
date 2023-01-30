package domain

import "time"

type LockRequest struct {
	Key string `json:"Key" valid:"required"`
	TTL time.Duration
}

type UnLockRequest struct {
	Key     string `json:"Key" valid:"required"`
	LockKey string
}

type LockResponse struct {
	LockKey string `json:"lockKey,omitempty"`
}
