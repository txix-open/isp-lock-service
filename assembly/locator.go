package assembly

import (
	"isp-lock-service/conf"
	"isp-lock-service/controller"
	"isp-lock-service/repository"
	"isp-lock-service/routes"
	"isp-lock-service/service"

	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/grpc/endpoint"
	"github.com/integration-system/isp-kit/log"
)

type Locator struct {
	cfg    conf.Remote
	logger log.Logger
}

func NewLocator(logger log.Logger, cfg conf.Remote) Locator {
	return Locator{
		cfg:    cfg,
		logger: logger,
	}
}

func (l Locator) Handler() *grpc.Mux {
	lockRepo := repository.NewLocker(l.logger, l.cfg)
	lockerService := service.NewLocker(l.logger, lockRepo)
	lockerController := controller.NewLocker(l.logger, lockerService)
	c := routes.Controllers{
		Locker: lockerController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.BodyLogger(l.logger))
	handler := routes.Handler(mapper, c)
	return handler
}
