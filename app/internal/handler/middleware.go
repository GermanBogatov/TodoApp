package handler

import (
	"errors"
	"fmt"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	userCtx = "userId"
)

type errorResponse struct {
	Message string `json:"message"'`
}

type statusResponse struct {
	Status string `json:"status"`
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}
	return strconv.Atoi(fmt.Sprintf("%s", id))
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func ValidateRequestSign(user model.UserDTO) error {
	fmt.Println("user:", user)
	if user.Id != 0 {
		return errors.New("Invalid id in request!")
	}
	if user.Username == "" {
		return errors.New("Invalid Username!")
	}
	if user.Name == "" {
		return errors.New("Invalid Name!")
	}
	if user.Password == "" {
		return errors.New("Invalid Password!")
	}
	if len(user.Password) < 6 {
		return errors.New("Invalid length Password!")
	}

	return nil
}
