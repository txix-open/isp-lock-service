package assembly

import (
	"context"

	"github.com/integration-system/isp-kit/observability/sentry"
	"isp-lock-service/conf"

	"github.com/integration-system/isp-kit/app"
	"github.com/integration-system/isp-kit/bootstrap"
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/log"
	"github.com/pkg/errors"
)

type Assembly struct {
	boot   *bootstrap.Bootstrap
	server *grpc.Server
	logger *log.Adapter
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
		a.boot.Fatal(errors.WithMessage(err, "upgrade remote config"))
	}

	a.logger.SetLevel(newCfg.LogLevel)

	logger := sentry.WrapErrorLogger(a.logger, a.boot.SentryHub)
	locator := NewLocator(logger, newCfg)
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
	}
}
