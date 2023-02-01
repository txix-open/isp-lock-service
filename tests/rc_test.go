package tests_test

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

// nolint: paralleltest
func TestOne(t *testing.T) {
	tst, required := test.New(t)
	rcli := NewRedis(tst)

	r, err := repository.NewRCWithClient("testPrefix", tst.Logger(), rcli)
	required.NoError(err)

	// success story
	key := time.Now().String()
	l, err := r.Lock(key, time.Second)
	required.NoError(err)

	_, err = r.UnLock(key, l)
	required.NoError(err)

	// look at wait
	_, err = r.Lock(key, time.Second)
	required.NoError(err)

	n := time.Now()
	_, err = r.Lock(key, time.Second)
	required.NoError(err)

	diff := time.Since(n)
	if diff < time.Second {
		t.Error("#2.3::error with second lock::", diff.String())
	}

	// look at error
	l, err = r.Lock(key, 10*time.Second)
	required.NoError(err)

	n = time.Now()
	_, err = r.Lock(key, time.Second)
	// diff = time.Since(n)

	required.Error(err)

	if err != nil {
		// nolint: errorlint
		if _, ok := err.(redsync.ErrTaken); !ok {
			t.Error("#3.3::", err)
		}
	}

	_, err = r.UnLock(key, l)
	required.NoError(err)
}
