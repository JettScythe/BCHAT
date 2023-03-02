package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func addUserRoutes() {
	users.GET("/:nameOrNumber", usersHandler)
}

func usersHandler(c *gin.Context) {
	nameOrNumber := c.Param("nameOrNumber")
	name, err := strconv.Atoi(nameOrNumber)
	if err != nil {
		// TODO: query the database to return a list of users
	}
	fmt.Println(name)
}
