package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func addUserRoutes() {
	users.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "users")
	})
	users.GET("/comments", func(c *gin.Context) {
		c.JSON(http.StatusOK, "users comments")
	})
	users.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, "users pictures")
	})
}
