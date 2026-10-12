package tests_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
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
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func TestOneLock(t *testing.T) {
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
	n := time.Now()
	_, err = r.Lock(ctx, key, 1)
	required.NoError(err)

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

func TestOneTryLock(t *testing.T) {
	t.Parallel()

	tst, required := test.New(t)
	rcli := NewRedis(tst)
	ctx := t.Context()

	r := repository.NewLocker(tst.Logger(), rcli, conf.Redis{Prefix: "testPrefix"}, conf.LockSettings{})

	// success story
	rnd, err := rand.Int(rand.Reader, big.NewInt(time.Now().Unix()))
	required.NoError(err)
	key := strconv.Itoa(int(rnd.Int64()))
	l, err := r.TryLock(ctx, key, 1)
	required.NoError(err)

	_, err = r.UnLock(ctx, key, l.LockKey)
	required.NoError(err)

	// second lock — no wait, immediate error on conflict
	resp, err := r.TryLock(ctx, key, 3)
	required.NoError(err)

	_, err = r.TryLock(ctx, key, 2)
	required.Error(err)
	if err != nil {
		required.Error(err, "fail lock")
	}

	_, err = r.UnLock(ctx, key, resp.LockKey)
	required.NoError(err)
}

func TestTryLockAfterTTL(t *testing.T) {
	t.Parallel()

	tst, required := test.New(t)
	rcli := NewRedis(tst)
	ctx := t.Context()

	r := repository.NewLocker(tst.Logger(), rcli, conf.Redis{Prefix: "testPrefix"}, conf.LockSettings{})

	rnd, err := rand.Int(rand.Reader, big.NewInt(time.Now().Unix()))
	required.NoError(err)
	key := strconv.Itoa(int(rnd.Int64()))
	_, err = r.TryLock(ctx, key, 1)
	required.NoError(err)

	time.Sleep(1200 * time.Millisecond)

	l2, err := r.TryLock(ctx, key, 1)
	required.NoError(err)

	_, err = r.UnLock(ctx, key, l2.LockKey)
	required.NoError(err)
}

func TestTryLockConcurrency(t *testing.T) {
	t.Parallel()
	tst, required := test.New(t)
	redis := NewRedis(tst)

	r := repository.NewLocker(tst.Logger(), redis, conf.Redis{Prefix: "testPrefix"}, conf.LockSettings{})

	var successCount atomic.Int64

	group := new(sync.WaitGroup)
	for range 10000 {
		group.Add(1)
		group.Go(func() {
			defer group.Done()
			resp, err := r.TryLock(t.Context(), "keyTryLock", 5)
			if err != nil {
				return
			}
			successCount.Add(1)
			_, err = r.UnLock(t.Context(), "keyTryLock", resp.LockKey)
			required.NoError(err)
		})
	}
	group.Wait()
	required.Positive(successCount.Load())
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
