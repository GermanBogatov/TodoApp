package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/TodoApp/app/internal/config"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var _ Helper = &helper{}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type RT struct {
	RefreshToken string `json:"refresh_token"`
}

type helper struct {
	Logger      logging.Logger
	clientRedis *redis.Client
}

func NewHelper(logger logging.Logger, client *redis.Client) Helper {
	return &helper{
		Logger:      logger,
		clientRedis: client,
	}
}

type Helper interface {
	GenerateAccessToken(u model.UserDTO) (string, string, error)
	UpdateRefreshToken(refreshToken string) (string, string, error)
}

func (h *helper) UpdateRefreshToken(refreshToken string) (string, string, error) {

	defer h.clientRedis.Del(context.Background(), refreshToken)

	userBytes := h.clientRedis.Get(context.Background(), refreshToken)
	fmt.Println("refresh: ", refreshToken)
	fmt.Println("userBytes: ", userBytes)
	var u model.UserDTO
	err := json.Unmarshal([]byte(userBytes.Val()), &u)
	if err != nil {
		return "", "", err
	}

	return h.GenerateAccessToken(u)

}

func (h *helper) GenerateAccessToken(u model.UserDTO) (string, string, error) {
	key := []byte(config.GetConfig().JWT.Secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(u.Id),
			Audience:  "users",
			ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
		},
		Username: u.Username,
	})

	accessToken, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}

	h.Logger.Info("create refresh token")
	refreshTokenUuid := uuid.New()
	userBytes, _ := json.Marshal(u)
	h.clientRedis.Set(context.Background(), refreshTokenUuid.String(), userBytes, 0)

	return accessToken, refreshTokenUuid.String(), err
}
