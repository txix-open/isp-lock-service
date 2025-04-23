package tests_test

import (
	"fmt"
	"testing"
	"time"

	"isp-lock-service/conf"
	"isp-lock-service/repository"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/txix-open/isp-kit/test"
)

func NewRedis(test *test.Test) *redis.Client {
	redisHost := test.Config().Optional().String("REDIS_HOST", "localhost")
	redisPort := test.Config().Optional().String("REDIS_PORT", "6379")
	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	return redis.NewClient(&redis.Options{Addr: addr})
}

func TestOne(t *testing.T) {
	t.Parallel()

	tst, required := test.New(t)
	rcli := NewRedis(tst)
	ctx := t.Context()

	r := repository.NewLocker(tst.Logger(), rcli, conf.Redis{Prefix: "testPrefix"}, conf.LockSettings{})

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

func TestConcurrency(t *testing.T) {
	t.Parallel()
	tst, required := test.New(t)
	redis := NewRedis(tst)

	r := repository.NewLocker(tst.Logger(), redis, conf.Redis{Prefix: "testPrefix"}, conf.LockSettings{})

	group, ctx := errgroup.WithContext(t.Context())
	group.SetLimit(32)
	for range 10000 {
		group.Go(func() error {
			resp, err := r.Lock(ctx, "key", 5)
			required.NoError(err)
			_, err = r.UnLock(ctx, "key", resp.LockKey)
			required.NoError(err)
			return nil
		})
	}
	err := group.Wait()
	required.NoError(err)
}
