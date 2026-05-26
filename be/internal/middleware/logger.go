
package middleware

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		methodColor := color.New(color.FgCyan).SprintFunc()
		statusColor := color.New(color.FgGreen).SprintFunc()
		pathColor := color.New(color.FgYellow).SprintFunc()
		timeColor := color.New(color.FgMagenta).SprintFunc()

		fmt.Printf(
			"%s %s %s %s\n",
			methodColor(c.Request.Method),
			pathColor(c.Request.URL.Path),
			statusColor(c.Writer.Status()),
			timeColor(duration),
		)
	}
}
