package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/serverapi/pkg/dto"
	"github.com/Todorov99/serverapi/pkg/entity"
	"github.com/Todorov99/serverapi/pkg/global"
	"github.com/Todorov99/serverapi/pkg/repository"
	"github.com/Todorov99/serverapi/pkg/server/config"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, registerDto dto.Register) error
	Login(ctx context.Context, loginDto dto.Login) (string, error)
}

type userService struct {
	logger         *logrus.Entry
	userRepository repository.UserRepository
}

func NewUserService() UserService {
	return &userService{
		logger:         logger.NewLogrus("userService", os.Stdout),
		userRepository: repository.NewUserRepository(),
	}
}

func (u *userService) Register(ctx context.Context, registerDto dto.Register) error {
	u.logger.Debugf("Registering user with username: %q", registerDto.UserName)
	_, err := u.userRepository.GetUserByUsername(ctx, registerDto.UserName)
	if errors.Is(err, global.ErrorUserWithUsernameNotExist) {
		passHash, err := getHash([]byte(registerDto.Password))
		if err != nil {
			return err
		}

		registerDto.Password = passHash

		userEntity := entity.User{}
		err = mapstructure.Decode(registerDto, &userEntity)
		if err != nil {
			return err
		}

		return u.userRepository.AddUser(ctx, userEntity)
	}

	return fmt.Errorf("user with username: %s already exists", registerDto.UserName)
}

func (u *userService) Login(ctx context.Context, loginDto dto.Login) (string, error) {
	user, err := u.userRepository.GetUserByUsername(ctx, loginDto.UserName)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	jwtCfg := config.GetJWTCfg()
	token, err := jwtCfg.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func getHash(pwd []byte) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(passHash), nil
}
