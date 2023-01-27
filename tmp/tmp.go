package main

import (
	"fmt"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func main() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	rs := redsync.New(pool)

	mutexname := "my-global-mutex"
	mutex := rs.NewMutex(mutexname, redsync.WithExpiry(time.Minute))

	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	var a string
	fmt.Scanf("%s", &a)

	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}
