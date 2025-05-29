package routes

import (
	"github.com/gaspartv/encurtador-de-url/internal/handlers"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1Users := v1.Group("/users")
	v1Users.GET("/", handlers.UsersHandler)

	router.Run("0.0.0.0:8080")
}
