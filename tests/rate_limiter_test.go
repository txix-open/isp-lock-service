package tests_test

import (
	"testing"
	"time"

	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/fake"
	"isp-lock-service/conf"
	"isp-lock-service/domain"
	"isp-lock-service/repository"
)

func TestRateLimiter(t *testing.T) {
	t.Parallel()
	tst, required := test.New(t)
	redis := NewRedis(tst)
	r := repository.NewRateLimiter(tst.Logger(), redis, conf.Remote{
		Redis: conf.Redis{Prefix: "test"},
		InMemLimiter: conf.InMemLimiter{
			ClearPeriodInSec:      10,
			LastUseThresholdInSec: 10,
		},
	})
	tst.T().Cleanup(r.Close)

	var (
		maxRps = 5
		key    = "key"
		ctx    = t.Context()
	)
	for i := range maxRps {
		resp, err := r.Limit(ctx, key, maxRps)
		required.NoError(err)
		exp := &domain.RateLimiterResponse{
			Allow:      true,
			Remaining:  maxRps - (i + 1),
			RetryAfter: -1,
		}
		required.EqualValues(exp, resp)
	}

	resp, err := r.Limit(ctx, key, maxRps)
	required.NoError(err)
	exp := &domain.RateLimiterResponse{
		Allow:      false,
		Remaining:  0,
		RetryAfter: resp.RetryAfter,
	}
	required.Greater(exp.RetryAfter, time.Duration(0))
	required.EqualValues(exp, resp)

	time.Sleep(time.Second)
	resp, err = r.Limit(ctx, key, maxRps)
	required.NoError(err)
	exp = &domain.RateLimiterResponse{
		Allow:      true,
		Remaining:  maxRps - 1,
		RetryAfter: -1,
	}
	required.EqualValues(exp, resp)

	for range 1000 {
		key := fake.It[string]()
		resp, err := r.Limit(ctx, key, 1)
		required.NoError(err)
		exp := &domain.RateLimiterResponse{
			Allow:      true,
			Remaining:  0,
			RetryAfter: -1,
		}
		required.EqualValues(exp, resp)
	}
}

func TestRateLimiterInMem(t *testing.T) {
	t.Parallel()
	tst, required := test.New(t)
	r := repository.NewRateLimiter(tst.Logger(), nil, conf.Remote{
		InMemLimiter: conf.InMemLimiter{
			ClearPeriodInSec:      10,
			LastUseThresholdInSec: 10,
		},
	})
	tst.T().Cleanup(r.Close)

	var (
		maxRps      = 5
		reqInterval = time.Second / time.Duration(maxRps)
		key         = "key"
		ctx         = t.Context()
	)
	for i := range 2 * maxRps {
		resp, err := r.LimitInMem(ctx, key, float64(maxRps), false)
		required.NoError(err)
		expPassAfter := reqInterval * time.Duration(i)
		required.True(resp.PassAfter <= expPassAfter)
	}

	for range 1000 {
		key := fake.It[string]()
		resp, err := r.LimitInMem(ctx, key, 1, false)
		required.NoError(err)
		required.True(resp.PassAfter <= 0)
	}
}

func TestRateLimiterInMemInfiniteKey(t *testing.T) {
	t.Parallel()
	tst, required := test.New(t)
	r := repository.NewRateLimiter(tst.Logger(), nil, conf.Remote{
		InMemLimiter: conf.InMemLimiter{
			ClearPeriodInSec:      2,
			LastUseThresholdInSec: 1,
		},
	})
	tst.T().Cleanup(r.Close)

	var (
		maxRps      = 1
		reqInterval = time.Second / time.Duration(maxRps)
		key         = "key"
		ctx         = t.Context()
	)
	for i := range 2 * maxRps {
		resp, err := r.LimitInMem(ctx, key, float64(maxRps), true)
		required.NoError(err)
		expPassAfter := reqInterval * time.Duration(i)
		required.True(resp.PassAfter <= expPassAfter)
	}

	for range 1000 {
		_, err := r.LimitInMem(ctx, key, 1, true)
		required.NoError(err)
	}
	<-time.After(time.Second * 4)

	resp, err := r.LimitInMem(ctx, key, 1, true)
	required.NoError(err)
	required.True(resp.PassAfter > 0)
}

func TestRateLimiterInMemLowRPSCase(t *testing.T) {
	t.Parallel()
	tst, required := test.New(t)
	r := repository.NewRateLimiter(tst.Logger(), nil, conf.Remote{
		InMemLimiter: conf.InMemLimiter{
			ClearPeriodInSec:      10,
			LastUseThresholdInSec: 10,
		},
	})
	tst.T().Cleanup(r.Close)

	var (
		maxRps      = 0.5
		reqInterval = time.Duration(float64(time.Second) / maxRps)
		key         = "key"
		ctx         = t.Context()
	)
	for i := range 3 {
		resp, err := r.LimitInMem(ctx, key, maxRps, false)
		required.NoError(err)
		expPassAfter := reqInterval * time.Duration(i+1)
		required.True(resp.PassAfter <= expPassAfter)
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	tst, _ := test.New(new(testing.T))
	redis := NewRedis(tst)
	r := repository.NewRateLimiter(tst.Logger(), redis, conf.Remote{
		Redis: conf.Redis{Prefix: "test"},
		InMemLimiter: conf.InMemLimiter{
			ClearPeriodInSec:      10,
			LastUseThresholdInSec: 10,
		},
	})
	b.Cleanup(r.Close)

	var (
		maxRps = 10
		key    = fake.It[string]()
		ctx    = b.Context()
	)
	b.ResetTimer()

	b.Run("Limit", func(b *testing.B) {
		for range b.N {
			_, err := r.Limit(ctx, key, maxRps)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})

	b.Run("LimitInMem", func(b *testing.B) {
		for range b.N {
			_, err := r.LimitInMem(ctx, key, float64(maxRps), false)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})
}

// nolint:gocognit
func BenchmarkRateLimiter2(b *testing.B) {
	tst, _ := test.New(new(testing.T))
	redis := NewRedis(tst)
	r := repository.NewRateLimiter(tst.Logger(), redis, conf.Remote{
		Redis: conf.Redis{Prefix: "test"},
		InMemLimiter: conf.InMemLimiter{
			ClearPeriodInSec:      10,
			LastUseThresholdInSec: 10,
		},
	})
	b.Cleanup(r.Close)

	var (
		keys = fake.It[[]string](fake.MinSliceSize(100), fake.MaxSliceSize(100))
		ctx  = b.Context()
	)
	b.ResetTimer()

	b.Run("Limit", func(b *testing.B) {
		for range b.N {
			for i := range 100 {
				maxRps := i + 1
				for _, key := range keys {
					_, err := r.Limit(ctx, key, maxRps)
					if err != nil {
						b.Fatalf("unexpected error for maxRps %d and key %s: %v", maxRps, key, err)
					}
				}
			}
		}
	})

	b.Run("LimitInMem", func(b *testing.B) {
		for range b.N {
			for i := range 100 {
				maxRps := i + 1
				for _, key := range keys {
					_, err := r.LimitInMem(ctx, key, float64(maxRps), false)
					if err != nil {
						b.Fatalf("unexpected error for maxRps %d and key %s: %v", maxRps, key, err)
					}
				}
			}
		}
	})
}
