package domain

import "time"

type Request struct {
	Key     string
	LockKey string
	TTL     time.Duration
}

type LockResponse struct {
	LockKey string `json:"lockKey,omitempty"`
}
