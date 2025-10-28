package domain

import (
	"time"
)

type RateLimiterRequest struct {
	Key    string `validate:"required"`
	MaxRps int    `validate:"required"`
}

type RateLimiterInMemRequest struct {
	Key    string  `validate:"required"`
	MaxRps float64 `validate:"required"`
}

type RateLimiterResponse struct {
	Allow      bool
	Remaining  int
	RetryAfter time.Duration
}

type RateLimiterInMemResponse struct {
	PassAfter time.Duration
}
