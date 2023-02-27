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
		c.JSON(http.StatusOK, utils.ParseRequest("cashid:bchchat.xyz/cashid?a=update&d=data&r=required&o=optional&x=nonce"))
	})
}
