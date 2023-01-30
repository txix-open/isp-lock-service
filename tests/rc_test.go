package tests

import (
	"fmt"
	"testing"
	"time"

	"isp-lock-service/rc"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/integration-system/isp-kit/test"
)

func NewRedis(test *test.Test) *redis.Client {
	redisHost := test.Config().Optional().String("REDIS_HOST", "localhost")
	redisPort := test.Config().Optional().String("REDIS_PORT", "6379")
	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	return redis.NewClient(&redis.Options{Addr: addr})
}

func TestOne(t *testing.T) {
	tst, _ := test.New(t)
	rcli := NewRedis(tst)

	r, err := rc.NewRCWithClient("testPrefix", time.Second*2, nil, rcli)
	if err != nil {
		t.Error(err)
	}

	// success story
	l, err := r.Lock("abc", time.Second)
	if err != nil {
		t.Error("#1.1::", err)
	}

	_, err = r.UnLock("abc", l)
	if err != nil {
		t.Error("#1.2::", err)
	}

	// look at wait
	l, err = r.Lock("abc", time.Second)
	if err != nil {
		t.Error("#2.1::", err)
	}

	n := time.Now()
	l, err = r.Lock("abc", time.Second)
	if err != nil {
		t.Error("#2.2::", err)
	}
	diff := time.Now().Sub(n)
	if diff < time.Second {
		t.Error("#2.3::error with second lock::", diff.String())
	}

	// look at error
	l, err = r.Lock("abc", time.Minute)
	if err != nil {
		t.Error("#3.1::", err)
	}

	n = time.Now()
	_, err = r.Lock("abc", time.Second)
	diff = time.Now().Sub(n)
	if err == nil {
		t.Error("#3.2::", err)
	} else {
		if _, ok := err.(redsync.ErrTaken); !ok {
			t.Error("#3.3::", err)
		}
	}

	_, err = r.UnLock("abc", l)
	if err != nil {
		t.Error("#3.4::", err)
	}

}
