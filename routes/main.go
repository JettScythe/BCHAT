package routes

import (
	"BCHat/database/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var router = gin.Default()
var users = router.Group("/users")
var demos = router.Group("/demos")

func signUp(ctx *gin.Context) {
	user := new(models.User)
	if err := ctx.Bind(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	models.Users = append(models.Users, user)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Signed up successfully.",
		"jwt": "123456789",
	})
}

func signIn(ctx *gin.Context) {
	user := new(models.User)
	if err := ctx.Bind(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	for _, u := range models.Users {
		if u.Username == user.Username && u.Password == user.Password {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "Signed in successfully.",
				"jwt": "123456789",
			})
			return
		}
	}
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Sign in failed."})
}

func GetRoutes() {
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	addUserRoutes()
	router.POST("signup", signUp)
	router.POST("signin", signIn)
	addDemoRoutes()
	err = router.Run()
	if err != nil {
		return
	}
}
