package service

import (
	"fmt"
	"time"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/entity"
	"github.com/Todorov99/server/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

//TODO get from .env file
var MySigningKey = []byte("unicorns")

type UserService interface {
	Register(registerDto dto.Register) error
	Login(loginDto dto.Login) (string, error)
}

type userService struct {
	userRepository repository.Repository
}

func NewUserService() UserService {
	return &userService{
		userRepository: repository.CreateUserRepository(),
	}
}

func (u *userService) Register(registerDto dto.Register) error {
	passHash, err := getHash([]byte(registerDto.Password))
	if err != nil {
		return err
	}

	registerDto.Password = passHash
	return u.userRepository.Add(registerDto.UserName, registerDto.Password, registerDto.FirstName, registerDto.LastName, registerDto.Email)
}

func (u *userService) Login(loginDto dto.Login) (string, error) {
	user, err := u.userRepository.GetByID(loginDto.UserName)
	if err != nil {
		return "", err
	}
	userEntity := entity.User{}
	err = mapstructure.Decode(user, &userEntity)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(loginDto.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	token, err := generateJWT(userEntity.ID)
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

func generateJWT(userID int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", fmt.Errorf("something Went Wrong: %w", err)
	}

	return tokenString, nil
}
