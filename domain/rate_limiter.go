package domain

import (
	"time"
)

type RateLimiterRequest struct {
	Key    string `validate:"required"`
	MaxRps int    `validate:"required"`
}

type RateLimiterResponse struct {
	Allow      bool
	Remaining  int
	RetryAfter time.Duration
}

type InMemRateLimiterResponse struct {
	PassAfter time.Duration
}
