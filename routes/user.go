package routes

import (
	"BCHChat/utils"
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
		c.JSON(http.StatusOK, utils.ParseRequest("cashid:bchchat.xyz/api/cashid?a=register&d=newsletter&r=i12p1c1&o=i458p3&x=95261230581"))
	})
}
