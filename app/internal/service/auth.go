package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/internal/storage"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	salt       = "sad342mslfd23412sdfsdf1234hgf"
	signingKey = ("HellowGerman! this is gin rest api")
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}
type AuthService struct {
	storage storage.Authorization
	logger  logging.Logger
}

func NewAuthService(storage storage.Authorization, logger logging.Logger) *AuthService {
	return &AuthService{
		storage: storage,
		logger:  logger,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.storage.CreateUser(ctx, user)
}

func (s *AuthService) GetUser(ctx context.Context, username, password string) (model.User, error) {
	return s.storage.GetUser(ctx, username, generatePasswordHash(password))
}
func (s *AuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	user, err := s.storage.GetUser(ctx, username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(ctx context.Context, accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token Claims are not of type")
	}

	return claims.UserId, nil

}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
