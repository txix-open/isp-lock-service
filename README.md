# Сервис распределенной блокировки

Общая идея состоит в том, чтобы реализовать поведение sync.Mutex на основе http вызова и через redis. Для этого используется библиотека [redsync](github.com/go-redsync/redsync). 

## Вызовы

Модуль имеет два вызова


### Блокировка 

- isp-lock-service/lock
- входящий запрос
  - key - строка, описывающая, что мы блокируем. ОБЯЗАТЕЛЬНОЕ
  - ttlInSec - число секунд после которых блокировка снимется автоматически. Если не передано, то будет использовано значение из конфига (.redis.defaultTimeOut). . ОБЯЗАТЕЛЬНОЕ
- ответ
  - lockKey - строка. ключ для разблокировки, чтобы тот, кто не ставил блокировку не мог ее снять 

### Пример вызова

```
POST {{proto}}://{{address}}/api/isp-lock-service/lock
Content-Type: application/json
X-Application-Token: {{token}}

{
    "key": "abc",
    "ttlInSec": 20
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
