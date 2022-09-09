package jwt

import (
	"github.com/GermanBogatov/TodoApp/app/internal/config"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

var _ Helper = &helper{}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type helper struct {
	Logger logging.Logger
	client *redis.Client
}

func NewHelper(logger logging.Logger, client *redis.Client) Helper {
	return &helper{
		Logger: logger,
		client: client,
	}
}

type Helper interface {
	GenerateAccessToken(u model.User) (string, error)
	UpdateRefreshToken() (string, error)
}

func (h *helper) UpdateRefreshToken() (string, error) {
	/*	defer h.RTCache.Del([]byte(rt.RefreshToken))

		userBytes, err := h.RTCache.Get([]byte(rt.RefreshToken))
		if err != nil {
			return nil, err
		}
		var u user_service.User
		err = json.Unmarshal(userBytes, &u)
		if err != nil {
			return nil, err
		}
		return h.GenerateAccessToken(u)*/
	return "", nil
}

func (h *helper) GenerateAccessToken(u model.User) (string, error) {
	key := []byte(config.GetConfig().JWT.Secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(u.Id),
			Audience:  "users",
			ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
		},
		Username: u.Username,
	})

	return token.SignedString(key)
}
