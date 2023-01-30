package rc

import (
	"context"
	"errors"
	"time"

	"isp-lock-service/conf"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/integration-system/isp-kit/log"
)

var ErrAlreadyLocked = errors.New("already locked")

type RC struct {
	cli     *redsync.Redsync
	timeOut time.Duration
	prefix  string
	l       *log.Adapter
}

func NewRCWithClient(cfg conf.Remote, l *log.Adapter, cli *goredislib.Client) (*RC, error) {
	r := RC{
		timeOut: cfg.Redis.DefaultTimeOut * time.Second,
		prefix:  cfg.Redis.Prefix,
		l:       l,
	}

	r.cli = redsync.New(goredis.NewPool(cli))

	return &r, nil
}

func NewRC(cfg conf.Remote, l *log.Adapter) (*RC, error) {
	r := RC{
		timeOut: cfg.Redis.DefaultTimeOut * time.Second,
		prefix:  cfg.Redis.Prefix,
		l:       l,
	}

	cli := goredislib.NewClient(&goredislib.Options{
		Addr:     cfg.Redis.Server,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if cfg.Redis.RedisSentinel != nil {
		cli = goredislib.NewFailoverClient(&goredislib.FailoverOptions{
			MasterName:       cfg.Redis.RedisSentinel.MasterName,
			SentinelAddrs:    cfg.Redis.RedisSentinel.Addresses,
			SentinelUsername: cfg.Redis.RedisSentinel.Username,
			SentinelPassword: cfg.Redis.RedisSentinel.Password,
			Password:         cfg.Redis.Password,
		})
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

	if ttl == 0 {
		ttl = rc.timeOut
	}

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

	return rc.cli.NewMutex(key, redsync.WithValue(lockKey)).Unlock()
}
