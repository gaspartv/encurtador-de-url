package handlers

import (
	"net/http"
	"os"
	"time"

	usersRepositories "github.com/gaspartv/encurtador-de-url/internal/repositories/users"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func LoginUsersHandlers(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := usersRepositories.FindUserRepository(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if req.Email != user.Email || req.Password != user.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": req.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"message": "login successful",
	})
}
