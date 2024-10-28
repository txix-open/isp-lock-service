package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"isp-lock-service/domain"
)

type dailyLimiterRepo interface {
	Increment(ctx context.Context, key string) (uint64, error)
	Set(ctx context.Context, key string, dailyLimit uint64) error
	GetLimit(ctx context.Context, key string) (uint64, error)
}

type dailyLimiter struct {
	repo dailyLimiterRepo
}

func NewDailyLimiter(repo dailyLimiterRepo) dailyLimiter { return dailyLimiter{repo: repo} }

func (s dailyLimiter) Increment(ctx context.Context, req domain.IncrementRequest) (*domain.IncrementResponse, error) {
	v, err := s.repo.Increment(ctx, s.makeKey(req.Key, req.Today))
	if err != nil {
		return nil, errors.WithMessage(err, "repo increment")
	}
	return &domain.IncrementResponse{Value: v}, nil
}

func (s dailyLimiter) Set(ctx context.Context, req domain.SetRequest) error {
	err := s.repo.Set(ctx, s.makeKey(req.Key, req.Today), req.Value)
	if err != nil {
		return errors.WithMessage(err, "repo set")
	}
	return nil
}

func (s dailyLimiter) GetLimit(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error) {
	v, err := s.repo.GetLimit(ctx, s.makeKey(req.Key, req.Today))
	if err != nil {
		return nil, errors.WithMessage(err, "repo get")
	}

	return &domain.GetResponse{Value: v}, nil
}

func (s dailyLimiter) makeKey(key string, today time.Time) string {
	return fmt.Sprintf("%s:%s", key, today.Format(time.DateOnly))
}
