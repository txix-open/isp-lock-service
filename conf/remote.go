package conf

import (
	"reflect"
	"time"

	"github.com/integration-system/isp-kit/log"
	"github.com/integration-system/isp-kit/rc/schema"
	"github.com/integration-system/jsonschema"
)

// nolint: gochecknoinits
func init() {
	schema.CustomGenerators.Register("logLevel", func(field reflect.StructField, t *jsonschema.Type) {
		t.Type = "string"
		t.Enum = []any{"debug", "info", "error", "fatal"}
	})
}

//	type Remote struct {
//		// Database dbx.Config
//		// Consumer Consumer
//		LogLevel log.Level `schemaGen:"logLevel" schema:"Уровень логирования"`
//	}
type Remote struct {
	LogLevel log.Level `schemaGen:"logLevel"  schema:"Уровень логирования"`
	Redis    struct {
		Server         string        `schemaGen:"server"  schema:"Адрес сервера redis"`
		Prefix         string        `schemaGen:"prefix"  schema:"Префикс ключей для модуля"`
		DefaultTimeOut time.Duration `schemaGen:"defaultTimeOut"  schema:"TTL по умолчанию, в секундах"`
	} `schemaGen:"redis"  schema:"Настройки redis"`
}

// type Consumer struct {
// 	Client grmqx.Connection
// 	Config grmqx.Consumer
// }
