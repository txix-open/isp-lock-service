package service

import (
	"context"

	"isp-lock-service/domain"
)

type limiterRepo interface {
	Limit(ctx context.Context, key string, maxRps int) (*domain.RateLimiterResponse, error)
	LimitInMem(ctx context.Context, key string, maxRps float64, infiniteKey bool) (*domain.RateLimiterInMemResponse, error)
}

type rateLimiter struct {
	repo limiterRepo
}

func NewRateLimiter(repo limiterRepo) rateLimiter {
	return rateLimiter{repo: repo}
}

func (s rateLimiter) Limit(ctx context.Context, req domain.RateLimiterRequest) (*domain.RateLimiterResponse, error) {
	return s.repo.Limit(ctx, req.Key, req.MaxRps) // nolint:wrapcheck
}

func (s rateLimiter) LimitInMem(ctx context.Context, req domain.RateLimiterInMemRequest) (*domain.RateLimiterInMemResponse, error) {
	return s.repo.LimitInMem(ctx, req.Key, req.MaxRps, req.InfiniteKey) // nolint:wrapcheck
}
