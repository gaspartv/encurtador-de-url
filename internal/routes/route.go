package routes

import (
	"github.com/gaspartv/encurtador-de-url/db"
	usersHandlers "github.com/gaspartv/encurtador-de-url/internal/handlers/users"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	db := db.SetupDatabase()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	v1 := router.Group("/api/v1")

	v1Users := v1.Group("/users")
	v1Users.GET("/list", usersHandlers.ListUsersHandler)
	v1Users.POST("/create", usersHandlers.CreateUsersHandler)

	v1Auth := v1.Group("/auth")
	v1Auth.POST("/sign-in", usersHandlers.LoginUsersHandlers)

	router.Run("0.0.0.0:8080")
}
