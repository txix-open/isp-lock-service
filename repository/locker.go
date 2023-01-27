package repository

import (
	"context"
	"fmt"

	"isp-lock-service/domain"
)

type Locker struct {
	// db db.DB
}

func (l Locker) Lock(ctx context.Context, req domain.Request) error {
	fmt.Println("call repo.Lock")
	return nil
}

func (l Locker) UnLock(ctx context.Context, req domain.Request) error {
	fmt.Println("call repo.UnLock")
	return nil
}

func NewLocker() Locker {
	return Locker{
		// db: db,
	}
}
