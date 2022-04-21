package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/entity"
	"github.com/Todorov99/sensorapi/pkg/global"
	"github.com/Todorov99/sensorapi/pkg/repository"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorapi/pkg/vault"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	userService := NewUserService()
	userCfg := config.GetUserCfg()

	vault, err := vault.New(config.GetVault())
	if err != nil {
		panic(err)
	}

	userSecret, err := vault.Get(userCfg.UserSecret)
	if err != nil {
		panic(err)
	}

	err = userService.Register(context.Background(), dto.Register{
		UserName:  userSecret.Name,
		Password:  userSecret.Value,
		FirstName: userCfg.FirstName,
		LastName:  userCfg.LastName,
		Email:     userCfg.Email,
	})

	if err != nil && strings.Contains(err.Error(), "already exists") {
		fmt.Println("Skipping user creation")
		return
	}

	if err != nil {
		panic(err)
	}

	err = userService.AddDevice(context.Background(), 1)
	if err != nil {
		panic(err)
	}
}

type UserService interface {
	Register(ctx context.Context, registerDto dto.Register) error
	Login(ctx context.Context, loginDto dto.Login) (string, error)
	AddDevice(ctx context.Context, userID int) error
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

func (d *userService) AddDevice(ctx context.Context, userID int) error {
	return d.userRepository.AddUserDvice(ctx, userID)
}

func getHash(pwd []byte) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(passHash), nil
}
