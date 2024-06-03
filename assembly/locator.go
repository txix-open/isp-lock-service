package assembly

import (
	goredislib "github.com/redis/go-redis/v9"
	"isp-lock-service/conf"
	"isp-lock-service/controller"
	"isp-lock-service/repository"
	"isp-lock-service/routes"
	"isp-lock-service/service"

	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"github.com/txix-open/isp-kit/log"
)

type Locator struct {
	redisCli goredislib.UniversalClient
	cfg      conf.Remote
	logger   log.Logger
}

func NewLocator(logger log.Logger, redisCli goredislib.UniversalClient, cfg conf.Remote) Locator {
	return Locator{
		redisCli: redisCli,
		cfg:      cfg,
		logger:   logger,
	}
}

func (l Locator) Handler() *grpc.Mux {
	lockRepo := repository.NewLocker(l.logger, l.redisCli, l.cfg.Redis)
	lockerService := service.NewLocker(l.logger, lockRepo)
	lockerController := controller.NewLocker(l.logger, lockerService)

	limiterRepo := repository.NewRateLimiter(l.logger, l.redisCli, l.cfg.Redis)
	limiterService := service.NewRateLimiter(limiterRepo)
	limiterController := controller.NewRateLimiter(limiterService)

	c := routes.Controllers{
		Locker:      lockerController,
		RateLimiter: limiterController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.BodyLogger(l.logger))
	handler := routes.Handler(mapper, c)
	return handler
}
