package handlers

import (
	"crypto/rand"
	"math/big"

	"github.com/gaspartv/encurtador-de-url/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required,url" description:"URL original a ser encurtada"`
}

func CreateUrlsHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var body CreateUrlRequest
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	userID := ctx.MustGet("userID").(float64)

	db := ctx.MustGet("db").(*gorm.DB)

	var user models.User
	query := db.Model(&models.User{})
	err := query.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(404, gin.H{"error": "Usuário não encontrado."})
			return
		}
		ctx.JSON(500, gin.H{"error": "internal server error."})
		return
	}

	randomUrl, err := GenerateRandomCode(4)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to generate random code"})
		return
	}

	url := models.Url{
		UserID:      int64(userID),
		OriginalUrl: body.OriginalUrl,
		ShortUrl:    randomUrl,
	}
	if err := db.Create(&url).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to create URL"})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "URL created successfully",
		"url": gin.H{
			"original_url": url.OriginalUrl,
			"short_url":    url.ShortUrl,
		},
	})
}

func GenerateRandomCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
