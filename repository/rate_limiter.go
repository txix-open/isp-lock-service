package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/pkg/errors"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/txix-open/isp-kit/log"
	"golang.org/x/time/rate"
	"isp-lock-service/conf"
	"isp-lock-service/domain"
)

type RateLimiter struct {
	logger log.Logger
	cli    *redis_rate.Limiter
	prefix string

	inMemLimiters map[string]*limiter
	mu            sync.Locker
	cancelFn      context.CancelFunc
}

type limiter struct {
	*rate.Limiter
	lastUse time.Time
}

func NewRateLimiter(logger log.Logger, cli goredislib.UniversalClient, cfg conf.Remote) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	limiter := &RateLimiter{
		logger:        logger,
		cli:           redis_rate.NewLimiter(cli),
		prefix:        cfg.Redis.Prefix + "::rate-limiter",
		inMemLimiters: make(map[string]*limiter),
		mu:            new(sync.Mutex),
		cancelFn:      cancel,
	}
	go limiter.clearInMemLimiters(ctx, cfg.InMemLimiter)

	return limiter
}

func (r *RateLimiter) Limit(ctx context.Context, key string, maxRps int) (*domain.RateLimiterResponse, error) {
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

func (r *RateLimiter) LimitInMem(_ context.Context, key string, maxRps int) (*domain.InMemRateLimiterResponse, error) {
	key = fmt.Sprintf("%s:::%d", key, maxRps)
	r.mu.Lock()
	defer r.mu.Unlock()

	lim, ok := r.inMemLimiters[key]
	if !ok {
		lim = &limiter{Limiter: rate.NewLimiter(rate.Limit(maxRps), 1)}
		r.inMemLimiters[key] = lim
	}
	lim.lastUse = time.Now()

	res := lim.ReserveN(lim.lastUse, 1)
	if !res.OK() {
		return nil, errors.Errorf("can't reserve time; key='%s'", key)
	}

	return &domain.InMemRateLimiterResponse{
		PassAfter: res.DelayFrom(lim.lastUse),
	}, nil
}

func (r *RateLimiter) clearInMemLimiters(ctx context.Context, cfg conf.InMemLimiter) {
	ctx = log.ToContext(ctx, log.String("worker", "limiterCleaner"))
	var (
		lastUseThreshold = time.Duration(cfg.LastUseThresholdInSec) * time.Second
		ticker           = time.NewTicker(time.Duration(cfg.ClearPeriodInSec) * time.Second)
	)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count := r.removeUnusedLimiters(lastUseThreshold)
			r.logger.Debug(ctx, fmt.Sprintf("removed %d unused limiters", count))
		case <-ctx.Done():
			r.logger.Info(ctx, errors.WithMessage(ctx.Err(), "limiter cleaner worker done"))
			return
		}
	}
}

func (r *RateLimiter) removeUnusedLimiters(lastUseThreshold time.Duration) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	limiterCount := len(r.inMemLimiters)
	for k, v := range r.inMemLimiters {
		if time.Since(v.lastUse) >= lastUseThreshold {
			delete(r.inMemLimiters, k)
		}
	}

	return limiterCount - len(r.inMemLimiters)
}

func (r *RateLimiter) Close() {
	r.cancelFn()
}
