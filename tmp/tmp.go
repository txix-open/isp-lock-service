package main

import (
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
	mutex := rs.NewMutex(mutexname, redsync.WithExpiry(time.Second))

	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	if ok, err := rs.NewMutex(mutexname, redsync.WithValue(mutex.Value())).Unlock(); !ok || err != nil {
		// if ok, err := mutex.Unlock(); !ok || err != nil {
		println("unlock failed#1")
	} else {
		println("ok#1")
	}

	// time.Sleep(time.Second * 2)
	//
	// if ok, err := rs.NewMutex(mutexname).Unlock(); !ok || err != nil {
	// 	println("unlock failed#2")
	// } else {
	// 	println("ok#2")
	// }
	//
	// if ok, err := mutex.Unlock(); !ok || err != nil {
	// 	println("unlock failed#3")
	// } else {
	// 	println("ok#3")
	// }

}
