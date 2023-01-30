package rc

import (
	"context"
	"fmt"
	"time"

	"isp-lock-service/conf"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/integration-system/isp-kit/log"
	"go.uber.org/zap"
)

type RC struct {
	cli     *redsync.Redsync
	timeOut time.Duration
	prefix  string
	l       *log.Adapter
}

func NewRCWithClient(prefix string, timeOut time.Duration, l *log.Adapter, cli *goredislib.Client) (*RC, error) {
	r := RC{
		timeOut: timeOut,
		prefix:  prefix,
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
		Username: cfg.Redis.UserName,
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
			Username:         cfg.Redis.UserName,
		})
	}

	r.cli = redsync.New(goredis.NewPool(cli))

	return &r, nil
}

func makeKey(prefix, key string) string {
	return prefix + "::" + key
}

func (rc *RC) log(level string, message interface{}, fields ...zap.Field) {
	if rc.l != nil {
		switch level {
		case "error":
			rc.l.Error(context.Background(), message, fields...)
		case "info":
			rc.l.Info(context.Background(), message, fields...)
		default:
			rc.l.Debug(context.Background(), message, fields...)
		}
	}
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
	rc.log("debug", "пробуем залочить "+key+" на "+ttl.String())

	mtx := rc.cli.NewMutex(key, redsync.WithExpiry(ttl))

	if err := mtx.Lock(); err != nil {
		rc.log("debug", fmt.Sprintf("err=%#v", err))
		return "", err
	}

	value := mtx.Value()
	rc.log("debug", "ключ для разблокировки "+value)
	return value, nil
}

// UnLock - функция снятия лока по ключу
//
//	key - суффикс ключа
//	lockKey - ключ, полученный в ответе из функции Lock
func (rc *RC) UnLock(key, lockKey string) (bool, error) {
	key = makeKey(rc.prefix, key)

	rc.log("debug", "пробуем разлочить "+key+"+"+lockKey)

	ok, err := rc.cli.NewMutex(key, redsync.WithValue(lockKey)).Unlock()

	rc.log("debug", fmt.Sprint("ok=", ok))
	if err != nil {
		rc.log("debug", fmt.Sprintf("err=%#v", err))
	}

	return ok, err
}
