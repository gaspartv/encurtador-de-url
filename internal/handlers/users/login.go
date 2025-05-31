package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gaspartv/encurtador-de-url/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	var user models.User

	db := ctx.MustGet("db").(*gorm.DB)
	query := db.Model(&models.User{})
	query = query.Where("email ILIKE ?", req.Email)
	query.Find(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if !strings.EqualFold(req.Email, user.Email) || err != nil {
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
