package repository

import (
	"context"

	"isp-lock-service/domain"

	"github.com/integration-system/isp-kit/log"
)

type Locker struct {
	// db db.DB
	logger log.Logger
}

func (l Locker) Lock(ctx context.Context, req domain.Request) error {
	l.logger.Debug(ctx, "call repo.Lock")
	return nil
}

func (l Locker) UnLock(ctx context.Context, req domain.Request) error {
	l.logger.Debug(ctx, "call repo.UnLock")
	return nil
}

func NewLocker(logger log.Logger) Locker {
	return Locker{
		// db: db,
		logger: logger,
	}
}
