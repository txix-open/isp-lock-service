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
		Server         string        `schemaGen:"server"  schema:"Адрес сервера redis, обязателен, если sentinel не указан"`
		Password       string        `schemaGen:"password"  schema:"Пароль для redis"`
		DB             int           `schemaGen:"db"  schema:"номер БД в redis"`
		Prefix         string        `schemaGen:"prefix"  schema:"Префикс ключей для модуля"`
		DefaultTimeOut time.Duration `schemaGen:"defaultTimeOut"  schema:"TTL по умолчанию, в секундах"`
		RedisSentinel  *struct {
			Addresses  []string `schema:"Адреса нод в кластере"`
			MasterName string   `schema:"Имя мастера"`
			Username   string   `schema:"Имя пользователя в sentinel"`
			Password   string   `schema:"Пароль в sentinel"`
		} `schema:"Настройки sentinel,обязательна, если redis.server не указан"`
	} `schemaGen:"redis"  schema:"Настройки redis"`
}

// type Consumer struct {
// 	Client grmqx.Connection
// 	Config grmqx.Consumer
// }
