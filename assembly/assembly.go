package assembly

import (
	"context"

	goredislib "github.com/redis/go-redis/v9"
	"github.com/txix-open/isp-kit/observability/sentry"
	"github.com/txix-open/isp-kit/rc"
	"isp-lock-service/conf"
	"isp-lock-service/repository"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/app"
	"github.com/txix-open/isp-kit/bootstrap"
	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/log"
)

type Assembly struct {
	boot            *bootstrap.Bootstrap
	server          *grpc.Server
	redisCli        *goredislib.Client
	logger          *log.Adapter
	rateLimiterRepo *repository.RateLimiter
}

func New(boot *bootstrap.Bootstrap) (*Assembly, error) {
	server := grpc.DefaultServer()
	return &Assembly{
		boot:            boot,
		server:          server,
		redisCli:        nil,
		logger:          boot.App.Logger(),
		rateLimiterRepo: nil,
	}, nil
}

func (a *Assembly) ReceiveConfig(ctx context.Context, remoteConfig []byte) error {
	if a.rateLimiterRepo != nil {
		a.rateLimiterRepo.Close()
	}

	newCfg, _, err := rc.Upgrade[conf.Remote](a.boot.RemoteConfig, remoteConfig)
	if err != nil {
		a.boot.Fatal(errors.WithMessage(err, "upgrade remote config"))
	}

	err = a.redisClient(newCfg.Redis)
	if err != nil {
		a.boot.Fatal(errors.WithMessage(err, "upgrade redis client"))
	}

	a.logger.SetLevel(newCfg.LogLevel)

	logger := sentry.WrapErrorLogger(a.logger, a.boot.SentryHub)
	locator := NewLocator(logger, a.redisCli, newCfg)
	locatorCfg := locator.Config()

	a.rateLimiterRepo = locatorCfg.rateLimiterRepo
	a.server.Upgrade(locatorCfg.handler)

	return nil
}

func (a *Assembly) Runners() []app.Runner {
	eventHandler := cluster.NewEventHandler().
		RemoteConfigReceiver(a)
	return []app.Runner{
		app.RunnerFunc(func(ctx context.Context) error {
			err := a.server.ListenAndServe(a.boot.BindingAddress)
			if err != nil {
				return errors.WithMessage(err, "listen ans serve grpc server")
			}
			return nil
		}),
		app.RunnerFunc(func(ctx context.Context) error {
			err := a.boot.ClusterCli.Run(ctx, eventHandler)
			if err != nil {
				return errors.WithMessage(err, "run cluster client")
			}
			return nil
		}),
	}
}

func (a *Assembly) Closers() []app.Closer {
	return []app.Closer{
		a.boot.ClusterCli,
		app.CloserFunc(func() error {
			a.server.Shutdown()
			return nil
		}),
		app.CloserFunc(func() error {
			if a.redisCli != nil {
				err := a.redisCli.Close()
				if err != nil {
					return errors.WithMessage(err, "shutdown redis client")
				}
			}
			return nil
		}),
		app.CloserFunc(func() error {
			if a.rateLimiterRepo != nil {
				a.rateLimiterRepo.Close()
			}
			return nil
		}),
	}
}

func (a *Assembly) redisClient(cfg conf.Redis) error {
	var (
		oldCli = a.redisCli
		newCli *goredislib.Client
	)
	if cfg.Sentinel != nil {
		newCli = goredislib.NewFailoverClient(&goredislib.FailoverOptions{
			MasterName:       cfg.Sentinel.MasterName,
			SentinelAddrs:    cfg.Sentinel.Addresses,
			SentinelUsername: cfg.Sentinel.Username,
			SentinelPassword: cfg.Sentinel.Password,
			Password:         cfg.Password,
			Username:         cfg.Username,
		})
	} else {
		newCli = goredislib.NewClient(&goredislib.Options{
			Addr:     cfg.Address,
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.Db,
		})
	}

	err := newCli.Ping(a.boot.App.Context()).Err()
	if err != nil {
		return errors.WithMessage(err, "ping new redis client")
	}
	a.redisCli = newCli

	if oldCli != nil {
		err := oldCli.Close()
		if err != nil {
			return errors.WithMessage(err, "close old redis client")
		}
	}
	return nil
}
