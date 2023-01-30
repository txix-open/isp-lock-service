package controller

import (
	"context"

	"isp-lock-service/domain"

	"github.com/integration-system/isp-kit/log"
)

type LockerService interface {
	Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error)
	UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error)
}

type Locker struct {
	s      LockerService
	logger log.Logger
}

// Lock
// @Tags locker
// @Summary выставляем лок на строку
// @Description `key` - строка для лока
// @Description `ttl` - число секунд после которых блокировка снимется автоматически
// @Description Возвращаем ключ для разблокировки
// @Accept json
// @Produce json
// @Success 200 {object} domain.LockResponse
// @Router /api/isp-lock-service/lock [POST]
func (l Locker) Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error) {
	l.logger.Debug(ctx, "call ctrl.Lock")
	return l.s.Lock(ctx, req)
}

// UnLock
// @Tags locker
// @Summary снимаем лок со строки
// @Description `key` - строка для лока
// @Description `locKKey` - ключ для разблокировки, полученный из Lock
// @Accept json
// @Produce json
// @Success 200
// @Router /api/isp-lock-service/unlock [POST]
func (l Locker) UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error) {
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
