package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/repository/query"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	AddUser(ctx context.Context, userEntity entity.User) error
	AddUserDvice(ctx context.Context, userID int) error
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserIDByUsername(ctx context.Context, name string) (int, error)
	GetUserIDByEmail(ctx context.Context, email string) (int, error)
}

type userRepository struct {
	logger        *logrus.Entry
	postgreClient *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		logger:        logger.NewLogrus("deviceRepository", os.Stdout),
		postgreClient: config.GetDatabaseCfg().GetPostgreClient(),
	}
}

func (d *userRepository) AddUserDvice(ctx context.Context, userID int) error {
	return executeModifyingQuery(ctx, query.UpdateUserDevice, d.postgreClient, userID)
}

func (u *userRepository) AddUser(ctx context.Context, userEntity entity.User) error {
	u.logger.Info("Adding user with username: %s", userEntity.UserName)
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

func (u *userRepository) GetUserIDByUsername(ctx context.Context, username string) (int, error) {
	return u.getUserID(ctx, query.GetUserIDByName, username)
}

func (u *userRepository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	return u.getUserID(ctx, query.GetUserIDByEmail, email)
}

func (u *userRepository) getUserID(ctx context.Context, query string, prop string) (int, error) {
	var userID int
	err := executeSelectQuery(ctx, query, u.postgreClient, &userID, prop)
	if err != nil {
		return 0, fmt.Errorf("failed getting user by: %q", prop)
	}
	if userID == 0 {
		return 0, global.ErrorUserWithUsernameNotExist
	}

	return userID, nil
}
