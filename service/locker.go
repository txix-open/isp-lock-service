package service

type LockerRepo interface {
	// All(ctx context.Context) ([]entity.Locker, error)
	// Get(ctx context.Context, id int) (*entity.Locker, error)
}

type Locker struct {
	repo LockerRepo
}

func NewLocker(repo LockerRepo) Locker {
	return Locker{
		repo: repo,
	}
}

// func (s Locker) All(ctx context.Context) ([]domain.Locker, error) {
// 	objects, err := s.repo.All(ctx)
// 	if err != nil {
// 		return nil, errors.WithMessage(err, "get all objects")
// 	}
// 	result := make([]domain.Locker, 0, len(objects))
// 	for _, object := range objects {
// 		d := domain.Locker{
// 			Name: object.Name,
// 		}
// 		result = append(result, d)
// 	}
// 	return result, nil
// }
//
// func (s Locker) Get(ctx context.Context, id int) (*domain.Locker, error) {
// 	object, err := s.repo.Get(ctx, id)
// 	if err != nil {
// 		return nil, errors.WithMessagef(err, "get object by id %d", id)
// 	}
// 	d := domain.Locker{
// 		Name: object.Name,
// 	}
// 	return &d, nil
// }
//
