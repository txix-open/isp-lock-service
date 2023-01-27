package repository

import (
	"context"
	"time"

	"isp-lock-service/domain"
	"isp-lock-service/rc"

	"github.com/integration-system/isp-kit/log"
)

type Locker struct {
	// db db.DB
	rc     *rc.RC
	logger log.Logger
}

func (l Locker) Lock(ctx context.Context, req domain.Request) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call repo.Lock")
	val, err := l.rc.Lock(req.Key, req.TTL*time.Second)
	if err != nil {
		return nil, err
	}
	return &domain.LockResponse{LockKey: val}, nil
}

func (l Locker) UnLock(ctx context.Context, req domain.Request) (domain.LockResponse, error) {
	l.logger.Debug(ctx, "call repo.UnLock")
	_, err := l.rc.UnLock(req.Key, req.LockKey)
	if err != nil {
		l.logger.Error(ctx, err)
		return domain.LockResponse{}, err
	}
	return domain.LockResponse{}, nil
}

func NewLocker(logger log.Logger, rc *rc.RC) Locker {
	return Locker{
		// db: db,
		rc:     rc,
		logger: logger,
	}
}
