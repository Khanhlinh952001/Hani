package auth

import "github.com/gin-gonic/gin"

const (
	userIDKey    = "authUserID"
	userEmailKey = "authUserEmail"
	userNameKey  = "authUserName"
	userPlanKey  = "authUserPlan"
	userGuestKey = "authGuest"
	guestIDKey   = "authGuestID"
	sessionIDKey = "authSessionID"
	claimsKey    = "authClaims"
)

func SetUser(c *gin.Context, claims *Claims) {
	c.Set(claimsKey, claims)
	c.Set(userIDKey, claims.UserID)
	c.Set(userEmailKey, claims.Email)
	c.Set(userNameKey, claims.Name)
	c.Set(userPlanKey, claims.Plan)
	c.Set(userGuestKey, claims.Guest)
	c.Set(guestIDKey, claims.GuestID)
	c.Set(sessionIDKey, claims.SessionID)
}

func ClaimsFromContext(c *gin.Context) (*Claims, bool) {
	v, ok := c.Get(claimsKey)
	if !ok {
		return nil, false
	}
	cl, ok := v.(*Claims)
	return cl, ok
}

func UserID(c *gin.Context) (int, bool) {
	v, ok := c.Get(userIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}

func IsGuest(c *gin.Context) bool {
	v, _ := c.Get(userGuestKey)
	g, _ := v.(bool)
	return g
}

func GuestID(c *gin.Context) string {
	v, _ := c.Get(guestIDKey)
	s, _ := v.(string)
	return s
}

func UserPlan(c *gin.Context) string {
	v, _ := c.Get(userPlanKey)
	s, _ := v.(string)
	return s
}

func UserEmail(c *gin.Context) string {
	v, _ := c.Get(userEmailKey)
	email, _ := v.(string)
	return email
}

func UserName(c *gin.Context) string {
	v, _ := c.Get(userNameKey)
	name, _ := v.(string)
	return name
}
