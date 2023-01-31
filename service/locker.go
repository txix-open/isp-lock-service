package service

import (
	"context"

	"isp-lock-service/domain"

	"github.com/integration-system/isp-kit/log"
)

type LockerRepo interface {
	Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error)
	UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error)
}

type Locker struct {
	repo   LockerRepo
	logger log.Logger
}

func (l Locker) Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call srv.Lock")
	return l.repo.Lock(ctx, req)
}

func (l Locker) UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call srv.UnLock")
	return l.repo.UnLock(ctx, req)
}

func NewLocker(logger log.Logger, repo LockerRepo) Locker {
	return Locker{
		repo:   repo,
		logger: logger,
	}
}
