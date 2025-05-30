package handlers

import (
	"net/http"

	usersRepositories "github.com/gaspartv/encurtador-de-url/internal/repositories/users"
	"github.com/gin-gonic/gin"
)

func ListUsersHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	userId := ctx.Query("id")
	email := ctx.Query("email")
	name := ctx.Query("name")
	disabled := ctx.Query("disabled")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	sortBy := ctx.Query("sort_by")
	sortOrder := ctx.Query("sort_order")

	users, err := usersRepositories.ListUsersRepository(
		userId,
		email,
		name,
		disabled,
		pageSize,
		pageNumber,
		sortBy,
		sortOrder,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
