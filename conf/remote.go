package conf

import (
	"reflect"

	"github.com/txix-open/isp-kit/log"
	"github.com/txix-open/isp-kit/rc/schema"
	"github.com/txix-open/jsonschema"
)

// nolint: gochecknoinits
func init() {
	schema.CustomGenerators.Register("logLevel", func(field reflect.StructField, s *jsonschema.Schema) {
		s.Type = "string"
		s.Enum = []any{"debug", "info", "error", "fatal"}
	})
}

type Remote struct {
	LogLevel log.Level `schemaGen:"logLevel" schema:"Уровень логирования"`
	Redis    Redis     `schema:"Настройки redis"`
}

type Redis struct {
	Address  string         `schema:"Адрес,обязательное, если sentinel не указан"`
	Username string         `schema:"Имя пользователя"`
	Password string         `schema:"Пароль"`
	Sentinel *RedisSentinel `schema:"Настройки sentinel,обязательно, если address не указан"`

	Db     int    `schema:"номер БД в redis"`
	Prefix string `schema:"Префикс ключей для модуля"`
}

type RedisSentinel struct {
	Addresses  []string `validate:"required" schema:"Адреса нод в кластере"`
	MasterName string   `validate:"required" schema:"Имя мастера"`
	Username   string   `schema:"Имя пользователя в sentinel"`
	Password   string   `schema:"Пароль в sentinel"`
}
