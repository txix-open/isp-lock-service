package domain

import "time"

type Request struct {
	Key     string `json:"Key" valid:"required"`
	LockKey string
	TTL     time.Duration
}

type LockResponse struct {
	LockKey string `json:"lockKey,omitempty"`
}
