package handler

import (
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/OzkrOssa/freeradius-api/internal/core/port"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	srv port.UserService
}

func NewUserHandler(srv port.UserService) *UserHandler {
	return &UserHandler{srv}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

func (uh *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := uh.srv.CreateUser(c, user)
	if err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userResponse := newUserResponse(user)
	response := newResponse(true, "user created successfully", userResponse)

	c.JSON(http.StatusOK, response)
}

type getUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	var req getUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := uh.srv.GetUser(c, req.ID)
	if err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userResponse := newUserResponse(user)

	response := newResponse(true, "user found successfully", userResponse)

	c.JSON(http.StatusOK, response)
}

type listUsersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

func (uh *UserHandler) ListUsers(c *gin.Context) {
	var req listUsersRequest
	var listUsers []userResponse

	if err := c.ShouldBindQuery(&req); err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	users, err := uh.srv.ListUsers(c, req.Skip, req.Limit)
	if err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	for _, user := range users {
		listUsers = append(listUsers, newUserResponse(&user))
	}

	total := uint64(len(users))
	meta := newMeta(total, req.Limit, req.Skip)

	data := map[string]interface{}{
		"meta":  meta,
		"users": listUsers,
	}

	response := newResponse(true, "success", data)

	c.JSON(http.StatusOK, response)
	return
}

type updateUserRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user := &domain.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	userUpdated, err := uh.srv.UpdateUser(c, user)
	if err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userResponse := newUserResponse(userUpdated)
	response := newResponse(true, "user updated successfully", userResponse)
	c.JSON(http.StatusOK, response)
	return
}

type deleteUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	var req deleteUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err := uh.srv.DeleteUser(c, req.ID); err != nil {
		m := parseError(err)
		errorResponse := newErrorResponse(m)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := newResponse(true, "user deleted successfully", nil)
	c.JSON(http.StatusOK, response)
	return
}
