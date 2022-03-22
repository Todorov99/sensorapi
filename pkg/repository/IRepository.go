package repository

import (
	"context"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
)

var repositoryLogger = logger.NewLogrus("repositroy", os.Stdout)

type IRepository interface {
	GetAll(ctx context.Context) (interface{}, error)
	GetByID(ctx context.Context, id int) (interface{}, error)
	Add(ctx context.Context, entity interface{}) error
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id int) error
}
