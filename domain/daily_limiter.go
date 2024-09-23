package domain

import (
	"time"
)

type IncrementRequest struct {
	Key   string    `validate:"required"`
	Today time.Time `validate:"required"`
}

type IncrementResponse struct {
	Value uint64
}

type SetRequest struct {
	Key   string    `validate:"required"`
	Value uint64    `validate:"required"`
	Today time.Time `validate:"required"`
}
