package service

import (
	"context"
)

type IService interface {
	GetAll(ctx context.Context) (interface{}, error)
	GetById(ctx context.Context, ID int) (interface{}, error)
	Add(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}) error
	Delete(ctx context.Context, ID int) (interface{}, error)
}
