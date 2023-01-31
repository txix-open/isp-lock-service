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

type Remote struct {
	LogLevel            log.Level     `schemaGen:"logLevel" schema:"Уровень логирования"`
	Redis               Redis         `schema:"Настройки redis"`
	DefaultTimeOutInSec time.Duration `schema:"TTLInSec по умолчанию, в секундах"`
}

type Redis struct {
	Address  string         `schema:"Адрес,обязательное, если sentinel не указан"`
	Username string         `schema:"Имя пользователя"`
	Password string         `schema:"Пароль"`
	Sentinel *RedisSentinel `schema:"Настройки sentinel,обязательно, если address не указан"`

	DB     int    `schema:"номер БД в redis"`
	Prefix string `schema:"Префикс ключей для модуля"`
}

type RedisSentinel struct {
	Addresses  []string `valid:"required" schema:"Адреса нод в кластере"`
	MasterName string   `valid:"required" schema:"Имя мастера"`
	Username   string   `schema:"Имя пользователя в sentinel"`
	Password   string   `schema:"Пароль в sentinel"`
}
