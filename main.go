package main

import (
	"isp-lock-service/assembly"
	"isp-lock-service/conf"
	"isp-lock-service/routes"

	"github.com/integration-system/isp-kit/bootstrap"
	"github.com/integration-system/isp-kit/shutdown"
)

var (
	version = "1.0.0"
)

// @title isp-lock-service
// @version 1.0.0
// @description Шаблон сервиса

// @license.name GNU GPL v3.0

// @host localhost:9000
// @BasePath /api/isp-lock-service

//go:generate swag init
//go:generate rm -f docs/swagger.json docs/docs.go
func main() {
	boot := bootstrap.New(version, conf.Remote{}, routes.EndpointDescriptors())
	app := boot.App
	logger := app.Logger()

	assembly, err := assembly.New(boot)
	if err != nil {
		boot.Fatal(err)
	}
	app.AddRunners(assembly.Runners()...)
	app.AddClosers(assembly.Closers()...)

	shutdown.On(func() {
		logger.Info(app.Context(), "starting shutdown")
		app.Shutdown()
		logger.Info(app.Context(), "shutdown completed")
	})

	err = app.Run()
	if err != nil {
		boot.Fatal(err)
	}
}
