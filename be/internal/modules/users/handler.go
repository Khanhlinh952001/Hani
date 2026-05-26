package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phone_number"`
	Provider    string `json:"provider"`
	ProviderId  string `json:"provider_id"`
	Avatar      string `json:"avatar"`
	Level       int    `json:"level"`
	Address     string `json:"address"`
	Status      int    `json:"status"`
	Role        int    `json:"role"`
}

func CreateUserHandler(c *gin.Context) {
	var req userRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Provider:    req.Provider,
		ProviderId:  req.ProviderId,
		Avatar:      req.Avatar,
		Level:       req.Level,
		Address:     req.Address,
		Status:      req.Status,
		Role:        req.Role,
	}

	if err := CreateUserService(user, req.Password); err != nil {
		if err.Error() == "email already taken" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func GetAllUsersHandler(c *gin.Context) {
	list, err := GetAllUsersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func GetUserByIDHandler(c *gin.Context) {
	user, err := GetUserByIDService(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

type updateUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Provider    string `json:"provider"`
	ProviderId  string `json:"provider_id"`
	Avatar      string `json:"avatar"`
	Level       int    `json:"level"`
	Address     string `json:"address"`
	Status      int    `json:"status"`
	Role        int    `json:"role"`
}

func UpdateUserHandler(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Provider:    req.Provider,
		ProviderId:  req.ProviderId,
		Avatar:      req.Avatar,
		Level:       req.Level,
		Address:     req.Address,
		Status:      req.Status,
		Role:        req.Role,
	}

	if err := UpdateUserService(c.Param("id"), user, req.Password); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	updated, _ := GetUserByIDService(c.Param("id"))
	c.JSON(http.StatusOK, updated)
}

func DeleteUserHandler(c *gin.Context) {
	if err := DeleteUserService(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
