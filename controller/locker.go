package controller

import (
	"context"

	"isp-lock-service/domain"

	"github.com/integration-system/isp-kit/log"
)

type LockerService interface {
	// All(ctx context.Context) ([]domain.Locker, error)
	// Get(ctx context.Context, id int) (*domain.Locker, error)
	Lock(ctx context.Context, req domain.Request) error
	UnLock(ctx context.Context, req domain.Request) error
}

type Locker struct {
	s      LockerService
	logger log.Logger
}

func (l Locker) Lock(ctx context.Context, req domain.Request) error {
	l.logger.Debug(ctx, "call ctrl.Lock")
	return l.s.Lock(ctx, req)
}

func (l Locker) UnLock(ctx context.Context, req domain.Request) error {
	l.logger.Debug(ctx, "call ctrl.UnLock")
	return l.s.UnLock(ctx, req)
}

func NewLocker(logger log.Logger, s LockerService) Locker {
	return Locker{
		s:      s,
		logger: logger,
	}
}

// All
// @Tags object
// @Summary Получить все объекты
// @Description Возвращает список объектов
// @Accept json
// @Produce json
// @Success 200 {array} domain.Locker
// @Failure 500 {object} domain.GrpcError
// @Router /object/all [POST]
// func (c Locker) All(ctx context.Context) ([]domain.Locker, error) {
// 	return c.s.All(ctx)
// }
