package auth

import "github.com/gin-gonic/gin"

const userIDKey = "authUserID"
const userEmailKey = "authUserEmail"
const userNameKey = "authUserName"

func SetUser(c *gin.Context, userID int, email, name string) {
	c.Set(userIDKey, userID)
	c.Set(userEmailKey, email)
	c.Set(userNameKey, name)
}

func UserID(c *gin.Context) (int, bool) {
	v, ok := c.Get(userIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
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
