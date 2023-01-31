package repository

import (
	"context"
	"fmt"
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
	rc     *RC
	logger log.Logger
}

func NewLocker(logger log.Logger, rc *RC) Locker {
	return Locker{
		rc:     rc,
		logger: logger,
	}
}

func (l Locker) Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call repo.Lock")
	val, err := l.rc.Lock(req.Key, time.Duration(req.TTLInSec)*time.Second)
	if err != nil {
		return nil, errors.WithMessage(err, "ошибка в lock")
	}
	return &domain.LockResponse{LockKey: val}, nil
}

func (l Locker) UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call repo.UnLock")
	_, err := l.rc.UnLock(req.Key, req.LockKey)
	if err != nil {
		return &domain.LockResponse{}, errors.WithMessage(err, "ошибка в unlock")
	}
	return &domain.LockResponse{}, nil
}

type RC struct {
	cli    *redsync.Redsync
	prefix string
	l      log.Logger
}

func NewRC(cfg conf.Remote, l log.Logger) (*RC, error) {
	r := RC{
		prefix: cfg.Redis.Prefix,
		l:      l,
	}

	cli := goredislib.NewClient(&goredislib.Options{
		Addr:     cfg.Redis.Address,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
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

	r.cli = redsync.New(goredis.NewPool(cli))

	return &r, nil
}

func NewRCWithClient(prefix string, l log.Logger, cli *goredislib.Client) (*RC, error) {
	r := RC{
		prefix: prefix,
		l:      l,
	}

	r.cli = redsync.New(goredis.NewPool(cli))

	return &r, nil
}

func makeKey(prefix, key string) string {
	return prefix + "::" + key
}

// Lock - функция установки лока по ключу
//
//	key - суффикс ключа
//	ttl - время жизни ключа
//
//	Возвращает ключ для разблокировки
func (rc *RC) Lock(key string, ttl time.Duration) (string, error) {
	key = makeKey(rc.prefix, key)

	rc.l.Debug(context.Background(), "пробуем залочить "+key+" на "+ttl.String())

	mtx := rc.cli.NewMutex(key, redsync.WithExpiry(ttl))

	if err := mtx.Lock(); err != nil {
		return "", err
	}

	value := mtx.Value()
	rc.l.Debug(context.Background(), "ключ для разблокировки "+value)

	return value, nil
}

// UnLock - функция снятия лока по ключу
//
//	key - суффикс ключа
//	lockKey - ключ, полученный в ответе из функции Lock
func (rc *RC) UnLock(key, lockKey string) (bool, error) {
	key = makeKey(rc.prefix, key)

	rc.l.Debug(context.Background(), "пробуем разлочить "+key+"+"+lockKey)

	ok, err := rc.cli.NewMutex(key, redsync.WithValue(lockKey)).Unlock()

	rc.l.Debug(context.Background(), fmt.Sprint("ok=", ok))

	return ok, err
}
