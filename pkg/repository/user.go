package repository

import (
	"database/sql"
	"fmt"

	"github.com/Todorov99/server/pkg/entity"
	"github.com/Todorov99/server/pkg/repository/query"
)

type userRepository struct {
	postgreClient *sql.DB
}

func (u *userRepository) GetAll() (interface{}, error) {
	return nil, nil
}

func (u *userRepository) GetByID(args ...string) (interface{}, error) {
	username := args[0]
	user := entity.User{}
	err := executeSelectQuery(query.GetUserByName, u.postgreClient, &user, username)
	if err != nil {
		return nil, fmt.Errorf("failed getting user with name: %q: %w", username, err)
	}

	return user, nil
}

func (u *userRepository) Add(args ...string) error {
	repositoryLogger.Info("Adding user...")
	username := args[0]

	userID := ""
	err := executeSelectQuery(query.GetUserIDByName, u.postgreClient, &userID, username)
	if err != nil {
		return fmt.Errorf("failed getting user with name: %q", username)
	}

	if userID != "" {
		return fmt.Errorf("user with name: %q already exists", username)
	}

	return executeModifyingQuery(query.InsertUser, u.postgreClient, args[0], args[1], args[2], args[3], args[4])
}

func (u *userRepository) Update(args ...string) error {
	return nil
}

func (u *userRepository) Delete(id string) (interface{}, error) {
	return nil, nil
}
