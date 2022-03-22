package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Todorov99/server/pkg/entity"
	"github.com/Todorov99/server/pkg/global"
	"github.com/Todorov99/server/pkg/repository/query"
	"github.com/Todorov99/server/pkg/server/config"
)

type UserRepository interface {
	AddUser(ctx context.Context, userEntity entity.User) error
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserIDByUsername(ctx context.Context, name string) (int64, error)
}

type userRepository struct {
	postgreClient *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}

func (u *userRepository) AddUser(ctx context.Context, userEntity entity.User) error {
	repositoryLogger.Info("Adding user with username: %s", userEntity.UserName)
	return executeModifyingQuery(ctx, query.InsertUser, u.postgreClient,
		userEntity.UserName, userEntity.Password,
		userEntity.FirstName, userEntity.LastName, userEntity.Email)
}

func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	user := entity.User{}
	err := executeSelectQuery(ctx, query.GetUserByName, u.postgreClient, &user, username)
	if err != nil {
		if errors.Is(err, global.ErrorObjectNotFound) {
			return user, global.ErrorUserWithUsernameNotExist
		}
		return user, err
	}
	return user, nil
}

func (u *userRepository) GetUserIDByUsername(ctx context.Context, username string) (int64, error) {
	var userID int64
	err := executeSelectQuery(ctx, query.GetUserIDByName, u.postgreClient, &userID, username)
	if err != nil {
		return 0, fmt.Errorf("failed getting user with name: %q", username)
	}

	if userID == 0 {
		return 0, global.ErrorUserWithUsernameNotExist
	}

	return userID, nil
}
