package handlers

import (
	"strings"

	"github.com/gaspartv/encurtador-de-url/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProfileUsersHandler(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

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

		var deletedAt interface{}
		if user.DeletedAt != nil {
			deletedAt = user.DeletedAt.Format("2006-01-02 15:04:05")
		} else {
			deletedAt = nil
		}

		ctx.JSON(200, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
			"deleted_at": deletedAt,
			"disabled":   user.Disabled,
		})
		return
	}
	ctx.JSON(401, gin.H{"error": "Authorization header missing or invalid"})
}
