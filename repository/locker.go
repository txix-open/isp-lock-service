package repository

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"isp-lock-service/conf"
	"isp-lock-service/domain"

	"github.com/pkg/errors"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/txix-open/isp-kit/log"
)

type Locker struct {
	logger log.Logger
	cli    *redsync.Redsync
	prefix string
}

func NewLocker(logger log.Logger, cli goredislib.UniversalClient, cfg conf.Redis) Locker {
	return Locker{
		logger: logger,
		cli:    redsync.New(goredis.NewPool(cli)),
		prefix: cfg.Prefix,
	}
}

func (l Locker) Lock(ctx context.Context, key string, ttl int) (*domain.LockResponse, error) {
	key = makeKey(l.prefix, key)

	mtx := l.cli.NewMutex(key, redsync.WithExpiry(time.Duration(ttl)*time.Second))

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
