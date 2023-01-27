package domain

import "time"

type Request struct {
	Key string
	TTL time.Duration
}

type LockResponse struct {
	LockKey string
}
