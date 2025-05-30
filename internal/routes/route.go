package routes

import (
	usersHandlers "github.com/gaspartv/encurtador-de-url/internal/handlers/users"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1Users := v1.Group("/users")
	v1Users.GET("/", usersHandlers.ListUsersHandler)

	v1Auth := v1.Group("/auth")
	v1Auth.POST("/sign-in", usersHandlers.LoginUsersHandlers)

	router.Run("0.0.0.0:8080")
}
