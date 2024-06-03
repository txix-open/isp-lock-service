package tests_test

import (
	"context"
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
	r := repository.NewRateLimiter(tst.Logger(), redis, conf.Redis{Prefix: "test"})

	var (
		maxRps = 5
		key    = "key"
	)
	for i := range maxRps {
		resp, err := r.Limit(context.Background(), key, maxRps)
		required.NoError(err)
		exp := &domain.RateLimiterResponse{
			Allow:      true,
			Remaining:  maxRps - (i + 1),
			RetryAfter: -1,
		}
		required.EqualValues(exp, resp)
	}

	resp, err := r.Limit(context.Background(), key, maxRps)
	required.NoError(err)
	exp := &domain.RateLimiterResponse{
		Allow:      false,
		Remaining:  0,
		RetryAfter: resp.RetryAfter,
	}
	required.Greater(exp.RetryAfter, time.Duration(0))
	required.EqualValues(exp, resp)

	time.Sleep(time.Second)
	resp, err = r.Limit(context.Background(), key, maxRps)
	required.NoError(err)
	exp = &domain.RateLimiterResponse{
		Allow:      true,
		Remaining:  maxRps - 1,
		RetryAfter: -1,
	}
	required.EqualValues(exp, resp)

	for range 1000 {
		key := fake.It[string]()
		resp, err := r.Limit(context.Background(), key, 1)
		required.NoError(err)
		exp := &domain.RateLimiterResponse{
			Allow:      true,
			Remaining:  0,
			RetryAfter: -1,
		}
		required.EqualValues(exp, resp)
	}
}
