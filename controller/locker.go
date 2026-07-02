package controller

import (
	"context"
	"errors"
	"net/http"

	"isp-lock-service/domain"

	"github.com/go-redsync/redsync/v4"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"github.com/txix-open/isp-kit/log"
)

type LockerService interface {
	Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error)
	TryLock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error)
	UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error)
}

type Locker struct {
	s      LockerService
	logger log.Logger
}

func NewLocker(logger log.Logger, s LockerService) Locker {
	return Locker{
		s:      s,
		logger: logger,
	}
}

// Lock
// @Tags locker
// @Summary выставляем лок на строку
// @Description Возвращаем ключ для разблокировки
// @Param key query string true "строка для лока"
// @Param ttlInSec query int true "число секунд после которых блокировка снимется автоматически"
// @Accept json
// @Produce json
// @Success 200 {object} domain.LockResponse
// @Router /api/isp-lock-service/lock [POST]
func (l Locker) Lock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error) {
	return l.s.Lock(ctx, req)
}

// TryLock
// @Tags locker
// @Summary однократная попытка взять лок без ретраев
// @Description Пытаемся взять лок один раз. В отличие от Lock, не делает повторных попыток. При успехе возвращаем ключ для разблокировки.
// @Param key query string true "строка для лока"
// @Param ttlInSec query int true "число секунд после которых блокировка снимется автоматически"
// @Accept json
// @Produce json
// @Success 200 {object} domain.LockResponse
// @Failure 423 {object} apierrors.Error
// @Router /api/isp-lock-service/try_lock [POST]
func (l Locker) TryLock(ctx context.Context, req domain.LockRequest) (*domain.LockResponse, error) {
	var errTaken *redsync.ErrTaken

	resp, err := l.s.TryLock(ctx, req)
	switch {
	case errors.As(err, &errTaken):
		return nil, apierrors.NewBusinessError(http.StatusLocked, err.Error(), err)
	case err != nil:
		return nil, err
	}

	return resp, nil
}

// UnLock
// @Tags locker
// @Summary снимаем лок со строки
// @Param key query string true "строка для лока"
// @Param lockKey query string true "ключ для разблокировки, полученный из Lock"
// @Accept json
// @Produce json
// @Success 200
// @Router /api/isp-lock-service/unlock [POST]
func (l Locker) UnLock(ctx context.Context, req domain.UnLockRequest) (*domain.LockResponse, error) {
	return l.s.UnLock(ctx, req)
}
