package handlers

import (
	"strings"

	"github.com/gaspartv/encurtador-de-url/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUsersHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var body map[string]interface{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	requiredFields := []string{"email", "password", "name"}
	for _, field := range requiredFields {
		if _, ok := body[field]; !ok {
			ctx.JSON(422, gin.H{"error": "missing field: " + field})
			return
		}
	}

	db := ctx.MustGet("db").(*gorm.DB)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Email:    strings.ToLower(body["email"].(string)),
		Password: string(hashedPassword),
		Name:     body["name"].(string),
	}
	if err := db.Create(&user).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(201, gin.H{"message": "user created successfully"})
}
