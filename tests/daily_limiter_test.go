package tests_test

import (
	"testing"

	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/fake"
	"isp-lock-service/conf"
	"isp-lock-service/domain"
	"isp-lock-service/repository"
	"isp-lock-service/service"
)

func TestDailyLimiter(t *testing.T) {
	t.Parallel()
	var (
		tst, required = test.New(t)
		redis         = NewRedis(tst)
		r             = repository.NewDailyLimiter(redis, conf.Redis{Prefix: "test"})
		svc           = service.NewDailyLimiter(r)
		ctx           = t.Context()
		req           = fake.It[domain.IncrementRequest]()
	)

	for i := range 1000 {
		resp, err := svc.Increment(ctx, req)
		required.NoError(err)
		required.EqualValues(i+1, resp.Value)
	}

	dailyLimits := fake.It[[]uint64](fake.MinSliceSize(1000))
	for _, v := range dailyLimits {
		err := svc.Set(ctx, domain.SetRequest{
			Key:   req.Key,
			Value: v,
			Today: req.Today,
		})
		required.NoError(err)

		resp, err := svc.Increment(ctx, req)
		required.NoError(err)
		required.EqualValues(v+1, resp.Value)
	}
}
