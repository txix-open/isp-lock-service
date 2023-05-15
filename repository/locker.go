package repository

import (
	"context"
	"time"

	"isp-lock-service/conf"
	"isp-lock-service/domain"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/integration-system/isp-kit/log"
	"github.com/pkg/errors"
)

type Locker struct {
	logger log.Logger
	cli    *redsync.Redsync
	prefix string
}

func NewLocker(logger log.Logger, cfg conf.Remote) Locker {
	cli := goredislib.NewClient(&goredislib.Options{
		Addr:     cfg.Redis.Address,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	if cfg.Redis.Sentinel != nil {
		cli = goredislib.NewFailoverClient(&goredislib.FailoverOptions{
			MasterName:       cfg.Redis.Sentinel.MasterName,
			SentinelAddrs:    cfg.Redis.Sentinel.Addresses,
			SentinelUsername: cfg.Redis.Sentinel.Username,
			SentinelPassword: cfg.Redis.Sentinel.Password,
			Password:         cfg.Redis.Password,
			Username:         cfg.Redis.Username,
		})
	}

	return Locker{
		logger: logger,
		cli:    redsync.New(goredis.NewPool(cli)),
		prefix: cfg.Redis.Prefix,
	}
}

func (l Locker) Lock(ctx context.Context, key string, ttl int) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call repo.Lock")

	key = makeKey(l.prefix, key)

	mtx := l.cli.NewMutex(key, redsync.WithExpiry(time.Duration(ttl)*time.Second))

	if err := mtx.Lock(); err != nil {
		return nil, errors.WithMessagef(err, "fail lock. key=%s", key)
	}

	value := mtx.Value()

	return &domain.LockResponse{LockKey: value}, nil
}

func (l Locker) UnLock(ctx context.Context, key, lockKey string) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call repo.UnLock")

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

func NewLockerWithClient(prefix string, l log.Logger, cli *goredislib.Client) (*Locker, error) {
	rc := redsync.New(goredis.NewPool(cli))

	return &Locker{
		logger: l,
		cli:    rc,
		prefix: prefix,
	}, nil
}
