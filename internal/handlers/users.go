package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UsersHandler(ctx *gin.Context) {
	userId := ctx.Query("id")

	ctx.JSON(http.StatusOK, gin.H{"userId": userId})
}
