package assembly

import (
	"isp-lock-service/conf"
	"isp-lock-service/controller"
	"isp-lock-service/repository"
	"isp-lock-service/routes"
	"isp-lock-service/service"

	goredislib "github.com/redis/go-redis/v9"

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

type locatorConfig struct {
	handler         *grpc.Mux
	rateLimiterRepo *repository.RateLimiter
}

func (l Locator) Config() locatorConfig {
	lockRepo := repository.NewLocker(l.logger, l.redisCli, l.cfg.Redis, l.cfg.LockSettings)
	lockerService := service.NewLocker(l.logger, lockRepo)
	lockerController := controller.NewLocker(l.logger, lockerService)

	rateLimiterRepo := repository.NewRateLimiter(l.logger, l.redisCli, l.cfg)
	rateLimiterService := service.NewRateLimiter(rateLimiterRepo)
	rateLimiterController := controller.NewRateLimiter(rateLimiterService)

	dailyLimiterRepo := repository.NewDailyLimiter(l.redisCli, l.cfg.Redis)
	dailyLimiterService := service.NewDailyLimiter(dailyLimiterRepo)
	dailyLimiterController := controller.NewDailyLimiter(dailyLimiterService)

	c := routes.Controllers{
		Locker:       lockerController,
		RateLimiter:  rateLimiterController,
		DailyLimiter: dailyLimiterController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.Log(l.logger, true))

	return locatorConfig{
		handler:         routes.Handler(mapper, c),
		rateLimiterRepo: rateLimiterRepo,
	}
}
