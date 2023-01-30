# Сервис распределенной блокировки

Общая идея состоит в том, чтобы реализовать поведение sync.Mutex на основе http вызова и через redis. Для этого используется библиотека [redsync](github.com/go-redsync/redsync). 

## Вызовы

Модуль имеет два вызова


### Блокировка 

- isp-lock-service/lock
- входящий запрос
  - key - строка, описывающая, что мы блокируем. ОБЯЗАТЕЛЬНОЕ
  - ttl - число секунд после которых блокировка снимется автоматически. Если не передано, то будет использовано значение из конфига (.redis.defaultTimeOut)
- ответ
  - lockKey - строка. ключ для разблокировки, чтобы тот, кто не ставил блокировку не мог ее снять 

### Пример вызова

```
POST {{proto}}://{{address}}/api/isp-lock-service/lock
Content-Type: application/json
X-Application-Token: {{token}}

{
    "key": "abc",
    "ttl": 20
}
```

### Пример ответа

```
{
    "lockKey": "Km1wmwaP3eKzqLbuEd46lQ=="
}
```

### Разблокировка

- isp-lock-service/unlock
- входящий запрос
    - key - строка, описывающая, что мы разблокируем. ОБЯЗАТЕЛЬНОЕ
    - lockKey - строка, ключ для разблокировки, полученный в ответ при блокировке. ОБЯЗАТЕЛЬНОЕ
- ответ
    - пустой

### Пример вызова

```
POST {{proto}}://{{address}}/api/isp-lock-service/unlock
Content-Type: application/json
X-Application-Token: {{token}}

{
    "key": "abc",
    "lockKey": "{{lockKey}}"
}
```

### Пример ответа

```
{
}
```

## Конфигурация

В конфигурации надо указать либо адрес для подключения к redis (.redis.server), либо настройки sentinel (.redis.sentinel.*)

```go
type Remote struct {
	LogLevel log.Level `schemaGen:"logLevel"  schema:"Уровень логирования"`
	Redis    struct {
		Server         string        `schemaGen:"server"  schema:"Адрес сервера redis, обязателен, если sentinel не указан"`
		UserName       string        `schemaGen:"userName"  schema:"Имя пользователя в  redis"`
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
```
