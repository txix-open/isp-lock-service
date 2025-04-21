package repository

import (
	"context"
	"time"

	"isp-lock-service/conf"
	"isp-lock-service/domain"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"golang.org/x/exp/rand"

	"github.com/pkg/errors"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/txix-open/isp-kit/log"
)

const (
	defaultMinLockRetryDelayMs = 50
	defaultMaxLockRetryDelayMs = 100
)

type Locker struct {
	logger              log.Logger
	cli                 *redsync.Redsync
	prefix              string
	minLockRetryDelayMs int
	maxLockRetryDelayMs int
}

func NewLocker(
	logger log.Logger,
	cli goredislib.UniversalClient,
	cfg conf.Redis,
	lockSettings conf.LockSettings,
) Locker {
	minLockRetryDelayMs := defaultMinLockRetryDelayMs
	if lockSettings.MinLockRetryDelayMs > 0 {
		minLockRetryDelayMs = lockSettings.MinLockRetryDelayMs
	}

	maxLockRetryDelayMs := defaultMaxLockRetryDelayMs
	if lockSettings.MaxLockRetryDelayMs > 0 {
		maxLockRetryDelayMs = lockSettings.MaxLockRetryDelayMs
	}

	return Locker{
		logger:              logger,
		cli:                 redsync.New(goredis.NewPool(cli)),
		prefix:              cfg.Prefix,
		minLockRetryDelayMs: minLockRetryDelayMs,
		maxLockRetryDelayMs: maxLockRetryDelayMs,
	}
}

func (l Locker) Lock(ctx context.Context, key string, ttl int) (*domain.LockResponse, error) {
	key = makeKey(l.prefix, key)

	mtx := l.cli.NewMutex(key,
		redsync.WithExpiry(time.Duration(ttl)*time.Second),
		redsync.WithRetryDelayFunc(
			func(_ int) time.Duration {
				delay := rand.Intn(l.maxLockRetryDelayMs-l.minLockRetryDelayMs) + l.minLockRetryDelayMs
				return time.Duration(delay) * time.Millisecond
			},
		),
	)

	if err := mtx.Lock(); err != nil {
		return nil, errors.WithMessagef(err, "fail lock. key=%s", key)
	}

	value := mtx.Value()

	return &domain.LockResponse{LockKey: value}, nil
}

func (l Locker) UnLock(ctx context.Context, key, lockKey string) (*domain.LockResponse, error) {
	key = makeKey(l.prefix, key)

	_, err := l.cli.NewMutex(key, redsync.WithValue(lockKey)).Unlock()
	if err != nil {
		return nil, errors.WithMessagef(err, "fail unlock. key=%s", key)
	}

	return &domain.LockResponse{}, nil
}

func makeKey(prefix, key string) string {
	return prefix + "::" + key
}
