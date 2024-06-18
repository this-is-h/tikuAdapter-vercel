// api/index.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vercel/go-bridge/go/bridge"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 定义你的路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin on Vercel!",
		})
	})

	// 处理请求
	bridge.ServeHTTP(router.ServeHTTP, w, r)
}
