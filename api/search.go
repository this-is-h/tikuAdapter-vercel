package handler

import (
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	c.String(http.StatusOK, "Hello from Gin!")
}
