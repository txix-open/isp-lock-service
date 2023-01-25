package controller

type LockerService interface {
	// All(ctx context.Context) ([]domain.Locker, error)
	// Get(ctx context.Context, id int) (*domain.Locker, error)
}

type Locker struct {
	s LockerService
}

func NewLocker(s LockerService) Locker {
	return Locker{
		s: s,
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
