package auth

import (
	"net/http"
	"strconv"

	"be/internal/billing"
	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func authResponse(c *gin.Context, status int, pair *TokenPair, user *users.User) {
	usage, _ := billing.GetUsageSnapshot(user.ID)
	c.JSON(status, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"expires_in":    pair.ExpiresIn,
		"token":         pair.Token,
		"user":          user,
		"usage":         usage,
	})
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
		Name:             req.Name,
		Email:            req.Email,
		Gender:           req.Gender,
		Provider:         "local",
		SubscriptionPlan: billing.PlanFree,
		IsActive:         true,
	}
	if err := users.CreateUserService(user, req.Password); err != nil {
		if err.Error() == "email already taken" {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pair, err := IssueTokensForUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	authResponse(c, http.StatusCreated, pair, user)
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
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account_banned"})
		return
	}

	pair, err := IssueTokensForUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	authResponse(c, http.StatusOK, pair, user)
}

func RefreshHandler(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, user, err := RefreshAccess(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if user != nil {
		authResponse(c, http.StatusOK, pair, user)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"expires_in":    pair.ExpiresIn,
		"token":         pair.Token,
	})
}

func LogoutHandler(c *gin.Context) {
	var req logoutRequest
	_ = c.ShouldBindJSON(&req)
	if req.RefreshToken != "" {
		_ = RevokeRefresh(req.RefreshToken)
	}
	if cl, ok := ClaimsFromContext(c); ok && cl.SessionID != "" {
		_ = RevokeSession(cl.SessionID)
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func MeHandler(c *gin.Context) {
	if IsGuest(c) {
		gid, err := uuid.Parse(GuestID(c))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guest"})
			return
		}
		usage, _ := billing.GetGuestUsageSnapshot(gid)
		c.JSON(http.StatusOK, gin.H{
			"guest": true,
			"plan":  billing.PlanGuest,
			"usage": usage,
		})
		return
	}

	userID, ok := UserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	usage, _ := billing.GetUsageSnapshot(userID)
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"usage": usage,
	})
}

type patchMeRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func PatchMeHandler(c *gin.Context) {
	userID, ok := UserID(c)
	if !ok || userID == 0 {
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
