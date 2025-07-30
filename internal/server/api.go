package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func healthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	})
}
