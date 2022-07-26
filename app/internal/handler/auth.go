package handler

import (
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body model.UserDTO true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input model.UserDTO
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := ValidateRequestSign(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.Service.Authorization.CreateUser(c.Request.Context(), input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body model.SignInDTO true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input model.SignInDTO
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.Service.GetUser(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, refreshToken, err := h.Helper.GenerateAccessToken(user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//h.helper.UpdateRefreshToken(refreshToken)

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) refresh(c *gin.Context) {
	refresh := c.Param("refresh_token")

	token, refreshToken, err := h.Helper.UpdateRefreshToken(refresh)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}
