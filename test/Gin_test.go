package test

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestGin(t *testing.T) {
	r := gin.Default()
	r.GET("/demo", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OOOOOK"})
	})
	r.Run(":8080")
}
