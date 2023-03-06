package routes

import (
	"BCHat/database/models"
	"BCHat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addDemoRoutes() {
	demos.POST("/signing-demo", func(c *gin.Context) {
		var req models.Payload
		if err := c.BindJSON(&req); err != nil {
			utils.InvalidateRequest(c, "REQUEST_BROKEN", "Invalid request payload")
		}
		_, err := utils.VerifySignature(req.Address, req.Request, req.Signature)
		if err != nil {
			utils.InvalidateRequest(c, "RESPONSE_INVALID_SIGNATURE", "Could not verify signature")
		}
		for key, element := range req.Metadata {
			fmt.Println("Key:", key, "=>", "Element:", element)
		}
		c.JSON(http.StatusOK, fmt.Sprintf("signing verified, metadata = %s", req.Metadata))

	})

	demos.GET("/parse-demo", func(c *gin.Context) {
		c.JSON(http.StatusOK, utils.ParseRequest("cashid:bchchat.xyz/api/cashid?a=register&d=newsletter&r=i12p1c1&o=i458p3&x=95261230581"))
	})

	demos.GET("/cashID/registerDemo", func(c *gin.Context) {
		metadata := map[string][]string{
			"optional": {"name", "family", "nickname", "age", "gender", "birthdate", "picture", "country", "state", "city", "street name", "street number", "residence", "coordinate", "email", "instant", "social", "phone", "postal"},
		}
		c.String(http.StatusOK, utils.CreateRequest(models.ServiceActions[3], "", metadata))
	})
}
