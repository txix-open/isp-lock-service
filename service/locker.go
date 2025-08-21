package service

import (
	"context"

	"isp-lock-service/domain"

	"github.com/txix-open/isp-kit/log"
)

type LockerRepo interface {
	Lock(ctx context.Context, key string, ttlInSec int) (*domain.LockResponse, error)
	UnLock(ctx context.Context, key, lockKey string) (*domain.LockResponse, error)
}

type Locker struct {
	repo   LockerRepo
	logger log.Logger
}

func NewLocker(logger log.Logger, repo LockerRepo) Locker {
	return Locker{
		repo:   repo,
		logger: logger,
	}
}

func (l Locker) Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error) {
	// nolint: wrapcheck
	return l.repo.Lock(ctx, req.Key, req.TTLInSec)
}

func (l Locker) UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error) {
	// nolint: wrapcheck
	return l.repo.UnLock(ctx, req.Key, req.LockKey)
}
