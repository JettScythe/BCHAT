package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func addUserRoutes() {
	users.GET("/:username", func(c *gin.Context) {
		c.JSON(http.StatusOK, "users")
	})
	users.GET("/:pagenum", func(c *gin.Context) {
		c.JSON(http.StatusOK, "users")
	})
}
