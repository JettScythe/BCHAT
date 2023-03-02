package routes

import (
	"github.com/gin-gonic/gin"
)

var router = gin.Default()
var users = router.Group("/users")
var demos = router.Group("/demos")

func GetRoutes() {
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	addUserRoutes()
	addDemoRoutes()
	err = router.Run()
	if err != nil {
		return
	}
}
