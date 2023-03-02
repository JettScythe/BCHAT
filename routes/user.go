package routes

import (
	"BCHat/database/models"
	"BCHat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addUserRoutes() {
	users.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "users")
	})

	users.POST("/register", func(c *gin.Context) {
		var req models.Payload
		if err := c.BindJSON(&req); err != nil {
			utils.InvalidateRequest(c, "REQUEST_BROKEN", "Invalid request payload")
		}
		_, err := utils.VerifySignature(req.Address, req.Data, req.Signature)
		if err != nil {
			utils.InvalidateRequest(c, "RESPONSE_INVALID_SIGNATURE", "Could not verify signature")
		}
	})

	users.GET("/parse-demo", func(c *gin.Context) {
		c.JSON(http.StatusOK, utils.ParseRequest("cashid:bchchat.xyz/api/cashid?a=register&d=newsletter&r=i12p1c1&o=i458p3&x=95261230581"))
	})

	users.GET("/cashID/demo", func(c *gin.Context) {
		metadata := map[string][]string{
			"required": {"name", "email"},
			"optional": {"age", "gender", "phone"},
		}
		c.String(http.StatusOK, utils.CreateRequest("login", "", metadata))
	})

}
