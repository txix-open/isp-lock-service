package controller

import (
	"context"

	"isp-lock-service/domain"
)

type limiterService interface {
	Limit(ctx context.Context, req domain.RateLimiterRequest) (*domain.RateLimiterResponse, error)
	LimitInMem(ctx context.Context, req domain.RateLimiterInMemRequest) (*domain.RateLimiterInMemResponse, error)
}

type RateLimiter struct {
	svc limiterService
}

func NewRateLimiter(svc limiterService) RateLimiter {
	return RateLimiter{svc: svc}
}

func (c RateLimiter) Limit(ctx context.Context, req domain.RateLimiterRequest) (*domain.RateLimiterResponse, error) {
	return c.svc.Limit(ctx, req)
}

func (c RateLimiter) LimitInMem(ctx context.Context, req domain.RateLimiterInMemRequest) (*domain.RateLimiterInMemResponse, error) {
	return c.svc.LimitInMem(ctx, req)
}
