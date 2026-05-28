package auth

import (
	"net/http"
	"strconv"

	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Gender   string `json:"gender" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if existing, _ := users.GetUserByEmailService(req.Email); existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	if !users.ValidGender(req.Gender) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gender must be male, female, or other"})
		return
	}

	user := &users.User{
		Name:     req.Name,
		Email:    req.Email,
		Gender:   req.Gender,
		Provider: "local",
	}
	if err := users.CreateUserService(user, req.Password); err != nil {
		if err.Error() == "email already taken" {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func LoginHandler(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.AuthenticateService(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func MeHandler(c *gin.Context) {
	userID, ok := UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

type patchMeRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func PatchMeHandler(c *gin.Context) {
	userID, ok := UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req patchMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" && req.Avatar == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
		return
	}

	patch := &users.User{Name: req.Name, Avatar: req.Avatar}
	if err := users.UpdateUserService(strconv.Itoa(userID), patch, ""); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
