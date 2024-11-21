package handler

import (
	"github.com/OzkrOssa/freeradius-api/internal/core/port"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	srv port.AuthService
}

func NewAuthHandler(srv port.AuthService) *AuthHandler {
	return &AuthHandler{srv}
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req loginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	token, err := ah.srv.Login(c, req.Email, req.Password)
	if err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := newAuthResponse(token)
	c.JSON(http.StatusOK, response)
}
