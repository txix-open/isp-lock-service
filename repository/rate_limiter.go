package repository

import (
	"context"

	"github.com/go-redis/redis_rate/v10"
	"github.com/pkg/errors"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/txix-open/isp-kit/log"
	"isp-lock-service/conf"
	"isp-lock-service/domain"
)

type rateLimiter struct {
	logger log.Logger
	cli    *redis_rate.Limiter
	prefix string
}

func NewRateLimiter(logger log.Logger, cli goredislib.UniversalClient, cfg conf.Redis) rateLimiter {
	return rateLimiter{
		logger: logger,
		cli:    redis_rate.NewLimiter(cli),
		prefix: cfg.Prefix + "::rate-limiter",
	}
}

func (r rateLimiter) Limit(ctx context.Context, key string, maxRps int) (*domain.RateLimiterResponse, error) {
	res, err := r.cli.Allow(ctx, makeKey(r.prefix, key), redis_rate.PerSecond(maxRps))
	if err != nil {
		return nil, errors.WithMessage(err, "allow request")
	}

	return &domain.RateLimiterResponse{
		Allow:      res.Allowed > 0,
		Remaining:  res.Remaining,
		RetryAfter: res.RetryAfter,
	}, nil
}
