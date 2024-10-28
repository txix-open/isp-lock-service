package controller

import (
	"context"

	"github.com/txix-open/isp-kit/grpc/apierrors"
	"isp-lock-service/domain"
)

type dailyLimiterService interface {
	Increment(ctx context.Context, req domain.IncrementRequest) (*domain.IncrementResponse, error)
	Set(ctx context.Context, req domain.SetRequest) error
	GetLimit(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error)
}

type DailyLimiter struct {
	svc dailyLimiterService
}

func NewDailyLimiter(svc dailyLimiterService) DailyLimiter { return DailyLimiter{svc: svc} }

func (c DailyLimiter) Increment(ctx context.Context, req domain.IncrementRequest) (*domain.IncrementResponse, error) {
	resp, err := c.svc.Increment(ctx, req)
	if err != nil {
		return nil, apierrors.NewInternalServiceError(err)
	}
	return resp, nil
}

func (c DailyLimiter) Set(ctx context.Context, req domain.SetRequest) error {
	err := c.svc.Set(ctx, req)
	if err != nil {
		return apierrors.NewInternalServiceError(err)
	}
	return nil
}

func (c DailyLimiter) Get(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error) {
	resp, err := c.svc.GetLimit(ctx, req)
	if err != nil {
		return nil, apierrors.NewInternalServiceError(err)
	}

	return resp, nil
}
