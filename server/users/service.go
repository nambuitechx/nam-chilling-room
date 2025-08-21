package users

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nambuitechx/nam-chilling-room-server/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (s *UserService) listUsers(username string, limit int, offset int) ([]User, error) {
	return s.UserRepository.selectUsers(username, limit, offset)
}

func (s *UserService) getUserByID(id string) (*User, error) {
	return s.UserRepository.selectUserByID(id)
}

func (s *UserService) createUser(username string, password string) (*User, error) {
	user, err := s.UserRepository.selectUserByUsername(username)

	if err != nil && err.Error() != "user not found" {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("username exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID: uuid.New().String(),
		Username: username,
		Password: string(hashedPassword),
	}

	return s.UserRepository.insertUser(newUser)
}

func (s *UserService) authenticate(username string, password string) (string, error) {
	user, err := s.UserRepository.selectUserByUsername(username)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	claimsStruct := utils.AuthorizedUserInfo {
		ID: user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "nam-chilling-room",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsStruct)

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// func (s *UserService) authorize(tokenString string) (*User, error) {
// 	claims, err := utils.ValidateTokenString(tokenString)

// 	if err != nil {
// 		return nil, err
// 	}

// 	user, err := s.UserRepository.SelectUserByUsername(claims.Username)

// 	return user, err
// }
