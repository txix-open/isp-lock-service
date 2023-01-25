package controller

import (
	"context"

	"isp-lock-service/domain"
	"isp-lock-service/entity"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ObjectService interface {
	All(ctx context.Context) ([]domain.Object, error)
	Get(ctx context.Context, id int) (*domain.Object, error)
}

type Object struct {
	s ObjectService
}

func NewObject(s ObjectService) Object {
	return Object{
		s: s,
	}
}

// All
// @Tags object
// @Summary Получить все объекты
// @Description Возвращает список объектов
// @Accept json
// @Produce json
// @Success 200 {array} domain.Object
// @Failure 500 {object} domain.GrpcError
// @Router /object/all [POST]
func (c Object) All(ctx context.Context) ([]domain.Object, error) {
	return c.s.All(ctx)
}

// GetById
// @Tags object
// @Summary Получить объект по его идентификатору
// @Description Возвращает объект
// @Accept json
// @Produce json
// @Param body body domain.ByIdRequest true "Идентификатор объекта"
// @Success 200 {object} domain.Object
// @Failure 404 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /object/get_by_id [POST]
func (c Object) GetById(ctx context.Context, req domain.ByIdRequest) (*domain.Object, error) {
	d, err := c.s.Get(ctx, req.Id)
	switch {
	case errors.Is(err, entity.ErrObjectNotFound):
		return nil, status.Errorf(codes.NotFound, "object by id '%d' not found", req.Id)
	default:
		return d, err
	}
}
