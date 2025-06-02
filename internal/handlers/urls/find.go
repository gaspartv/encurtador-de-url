package handlers

import (
	"net/http"

	"github.com/gaspartv/encurtador-de-url/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindUrlsHandler(ctx *gin.Context) {
	shortUrl := ctx.Param("short")

	var url models.Url

	db := ctx.MustGet("db").(*gorm.DB)

	if err := db.Model(&models.Url{}).Where("short_url = ?", shortUrl).First(&url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	url = models.Url{
		ID:          url.ID,
		UserID:      url.UserID,
		OriginalUrl: url.OriginalUrl,
		ShortUrl:    url.ShortUrl,
		Clicks:      url.Clicks + 1,
		CreatedAt:   url.CreatedAt,
	}
	if err := db.Model(&models.Url{}).Where("id = ?", url.ID).Update("clicks", url.Clicks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update click count"})
		return
	}

	ctx.Redirect(http.StatusFound, url.OriginalUrl)
}
