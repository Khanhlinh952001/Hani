package billing

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PlansHandler(c *gin.Context) {
	var plans []PlanLimit
	if err := dbListPlans(&plans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func ResetUsageHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := ResetUserUsage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	snap, _ := GetUsageSnapshot(id)
	c.JSON(http.StatusOK, snap)
}
