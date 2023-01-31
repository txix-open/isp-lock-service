package tests

import (
	"fmt"
	"testing"
	"time"

	"isp-lock-service/repository"

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
	tst, required := test.New(t)
	rcli := NewRedis(tst)

	r, err := repository.NewRCWithClient("testPrefix", tst.Logger(), rcli)
	required.NoError(err)

	// success story
	l, err := r.Lock("abc", time.Second)
	required.NoError(err)

	_, err = r.UnLock("abc", l)
	required.NoError(err)

	// look at wait
	l, err = r.Lock("abc", time.Second)
	required.NoError(err)

	n := time.Now()
	l, err = r.Lock("abc", time.Second)
	required.NoError(err)

	diff := time.Now().Sub(n)
	if diff < time.Second {
		t.Error("#2.3::error with second lock::", diff.String())
	}

	// look at error
	l, err = r.Lock("abc", 10*time.Second)
	required.NoError(err)

	n = time.Now()
	_, err = r.Lock("abc", time.Second)
	diff = time.Now().Sub(n)

	required.Error(err)

	if err != nil {
		if _, ok := err.(redsync.ErrTaken); !ok {
			t.Error("#3.3::", err)
		}
	}

	_, err = r.UnLock("abc", l)
	required.NoError(err)

}
