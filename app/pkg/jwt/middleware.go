package jwt

import (
	"errors"
	"github.com/GermanBogatov/TodoApp/app/internal/config"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logging.GetLogger()
		authHeader := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logger.Error("Malformed token")
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("malformed token"))
			return
		}
		logger.Debug("create jwt verifier")
		accessToken := authHeader[1]
		key := []byte(config.GetConfig().JWT.Secret)

		token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return key, nil
		})
		if err != nil {
			unauthorized(c.Writer, err)
			return
		}

		if !token.Valid {
			logger.Error("token has been inspired")
			unauthorized(c.Writer, err)
			return
		}
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			unauthorized(c.Writer, err)
			return
		}
		c.Set("userId", claims.Id)
		c.Next()

	}
}

func unauthorized(w http.ResponseWriter, err error) {
	logging.GetLogger().Error(err)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}
