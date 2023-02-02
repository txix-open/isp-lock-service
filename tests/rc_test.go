package tests_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"isp-lock-service/repository"

	"github.com/go-redis/redis/v8"
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
	ctx := context.Background()

	r, err := repository.NewLockerWithClient("testPrefix", tst.Logger(), rcli)
	required.NoError(err)

	// success story
	key := time.Now().String()
	l, err := r.Lock(ctx, key, 1)
	required.NoError(err)

	_, err = r.UnLock(ctx, key, l.LockKey)
	required.NoError(err)

	// look at wait
	_, err = r.Lock(ctx, key, 1)
	required.NoError(err)

	n := time.Now()
	_, err = r.Lock(ctx, key, 1)
	required.NoError(err)

	diff := time.Since(n)
	required.Greater(diff, time.Second)

	// look at error
	l, err = r.Lock(ctx, key, 10)
	required.NoError(err)

	_, err = r.Lock(ctx, key, 1)

	required.Error(err)

	if err != nil {
		required.Error(err, "fail lock")
	}

	_, err = r.UnLock(ctx, key, l.LockKey)
	required.NoError(err)
}
