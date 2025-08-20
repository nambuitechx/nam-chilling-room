package users

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (s *UserService) ListUsers(username string, limit int, offset int) ([]User, error) {
	return s.UserRepository.SelectUsers(username, limit, offset)
}

func (s *UserService) GetUserByID(id string) (*User, error) {
	return s.UserRepository.SelectUserByID(id)
}

func (s *UserService) CreateUser(username string, password string) (*User, error) {
	user, err := s.UserRepository.SelectUserByUsername(username)

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

	return s.UserRepository.InsertUser(newUser)
}

func (s *UserService) Authenticate(username string, password string) (string, error) {
	user, err := s.UserRepository.SelectUserByUsername(username)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	claimsStruct := AuthorizedUserInfo {
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

func (s *UserService) Authorize(tokenString string) (*User, error) {
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(AuthorizedUserInfo)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	user, err := s.UserRepository.SelectUserByUsername(claims.Username)

	return user, err
}
