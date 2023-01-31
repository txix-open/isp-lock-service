package assembly

import (
	"context"

	"isp-lock-service/conf"
	"isp-lock-service/repository"

	"github.com/integration-system/isp-kit/app"
	"github.com/integration-system/isp-kit/bootstrap"
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/log"
	"github.com/pkg/errors"
)

type Assembly struct {
	boot     *bootstrap.Bootstrap
	server   *grpc.Server
	logger   *log.Adapter
	redisCli *repository.RC
}

func New(boot *bootstrap.Bootstrap) (*Assembly, error) {
	server := grpc.DefaultServer()
	return &Assembly{
		boot:   boot,
		server: server,
		logger: boot.App.Logger(),
	}, nil
}

func (a *Assembly) ReceiveConfig(ctx context.Context, remoteConfig []byte) error {
	var (
		newCfg  conf.Remote
		prevCfg conf.Remote
	)
	err := a.boot.RemoteConfig.Upgrade(remoteConfig, &newCfg, &prevCfg)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "upgrade remote config"))
	}

	a.logger.SetLevel(newCfg.LogLevel)

	a.redisCli, err = repository.NewRC(newCfg, a.logger)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "error on connect to redis"))
	}

	locator := NewLocator(a.logger, a.redisCli)
	handler := locator.Handler()
	a.server.Upgrade(handler)

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
			return nil
		}),
	}
}
