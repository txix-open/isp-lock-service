package assembly

import (
	"isp-lock-service/controller"
	"isp-lock-service/repository"
	"isp-lock-service/routes"
	"isp-lock-service/service"

	"github.com/integration-system/isp-kit/db"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/grpc/endpoint"
	"github.com/integration-system/isp-kit/log"
)

type DB interface {
	db.DB
	db.Transactional
}

type Locator struct {
	// db     DB
	logger log.Logger
}

func NewLocator(logger log.Logger) Locator {
	return Locator{
		// db:     db,
		logger: logger,
	}
}

func (l Locator) Handler() *grpc.Mux {
	// objectRepo := repository.NewObject()
	lockRepo := repository.NewLocker(l.logger)
	// objectService := service.NewObject(objectRepo)
	lockerService := service.NewLocker(l.logger, lockRepo)
	// objectController := controller.NewObject(objectService)
	lockerController := controller.NewLocker(l.logger, lockerService)
	c := routes.Controllers{
		// Object: objectController,
		Locker: lockerController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.BodyLogger(l.logger))
	handler := routes.Handler(mapper, c)
	return handler
}

// func (l Locator) (consumerCfg conf.Consumer) grmqx.Config {
// 	// txManager := transaction.NewManager(l.db)
// 	// msgService := service.NewMessage(l.logger, txManager)
// 	// msgController := controller.NewMessage(msgService)
//
// 	// handler := grmqx.NewResultHandler(l.logger, msgController)
// 	// consumer := consumerCfg.Config.DefaultConsumer(handler, grmqx.ConsumerLog(l.logger))
//
// 	// brokerConfig := grmqx.NewConfig(
// 	// consumerCfg.Client.Url(),
// 	// grmqx.WithConsumers(consumer),
// 	// grmqx.WithDeclarations(grmqx.TopologyFromConsumers(consumerCfg.Config)),
// 	// )
// 	// return brokerConfig
// 	return grmqx.Config{}
// }
