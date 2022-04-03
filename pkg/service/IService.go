package service

import (
	"context"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
)

var serviceLogger = logger.NewLogrus("service", os.Stdout)

type IService interface {
	GetAll(ctx context.Context) (interface{}, error)
	GetById(ctx context.Context, ID int) (interface{}, error)
	Add(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}) error
	Delete(ctx context.Context, ID int) (interface{}, error)
}
