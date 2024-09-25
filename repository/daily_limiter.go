package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"isp-lock-service/conf"
)

const (
	dailyLimiterKeyTtl = 24 * time.Hour
)

type dailyLimiter struct {
	redisCli  redis.UniversalClient
	keyPrefix string
}

func NewDailyLimiter(redisCli redis.UniversalClient, cfg conf.Redis) dailyLimiter {
	return dailyLimiter{
		redisCli:  redisCli,
		keyPrefix: cfg.Prefix + "::daily-limit",
	}
}

func (r dailyLimiter) Increment(ctx context.Context, key string) (uint64, error) {
	key = makeKey(r.keyPrefix, key)
	v, err := r.redisCli.Incr(ctx, key).Uint64()
	if err != nil {
		return 0, errors.WithMessage(err, "redis cli incr")
	}

	if v == 1 {
		err := r.redisCli.ExpireNX(ctx, key, dailyLimiterKeyTtl).Err()
		if err != nil {
			return 0, errors.WithMessage(err, "redis cli expire nx")
		}
	}

	return v, nil
}

func (r dailyLimiter) Set(ctx context.Context, key string, dailyLimit uint64) error {
	err := r.redisCli.Set(ctx, makeKey(r.keyPrefix, key), dailyLimit, dailyLimiterKeyTtl).Err()
	if err != nil {
		return errors.WithMessage(err, "redis cli set nx")
	}
	return nil
}
