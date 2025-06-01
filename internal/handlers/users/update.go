package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gaspartv/encurtador-de-url/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UpdateUserRequest struct {
	ID       uint    `json:"id" binding:"required" description:"ID do usuário"`
	Email    string  `json:"email" binding:"required,email" description:"E-mail do usuário"`
	Name     string  `json:"name" binding:"required" description:"Nome do usuário"`
	Password *string `json:"password,omitempty" description:"Senha do usuário (opcional)"`
	Deleted  *bool   `json:"deleted" description:"Se o usuário está deletado"`
	Disabled bool    `json:"disabled" description:"Se o usuário está desabilitado"`
}

func UpdateUsersHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var body UpdateUserRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			ctx.JSON(400, gin.H{
				"error": "JSON mal formatado.",
			})
			return

		case errors.As(err, &unmarshalTypeError):
			ctx.JSON(422, gin.H{
				"error": fmt.Sprintf("Campo '%s' está com tipo inválido.", unmarshalTypeError.Field),
			})
			return

		default:
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	db := ctx.MustGet("db").(*gorm.DB)

	var userFound models.User
	query := db.Model(&models.User{})
	err := query.Where("id = ?", body.ID).First(&userFound).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(404, gin.H{"error": "Usuário não encontrado."})
			return
		}
		ctx.JSON(500, gin.H{"error": "internal server error."})
		return
	}

	if body.Email != userFound.Email {
		var existingUser models.User
		err := db.Model(&models.User{}).Where("email = ?", body.Email).First(&existingUser).Error
		if err == nil {
			ctx.JSON(409, gin.H{"error": "E-mail já está em uso."})
			return
		}
		if err != gorm.ErrRecordNotFound {
			ctx.JSON(500, gin.H{"error": "internal server error."})
			return
		}
	}

	if body.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Falha ao gerar a senha."})
			return
		}
		userFound.Password = string(hashedPassword)
	}

	if body.Name == "" {
		ctx.JSON(422, gin.H{"error": "Nome não pode ser vazio."})
		return
	}

	if body.Deleted != nil {
		if *body.Deleted {
			now := time.Now()
			userFound.DeletedAt = &now
		} else {
			userFound.DeletedAt = nil
		}
	}

	user := models.User{
		ID:        userFound.ID,
		CreatedAt: userFound.CreatedAt,
		UpdatedAt: userFound.UpdatedAt,
		Email:     body.Email,
		Password:  userFound.Password,
		Name:      body.Name,
		DeletedAt: userFound.DeletedAt,
		Disabled:  body.Disabled,
	}
	if err := db.Save(&user).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Falha ao atualizar o usuário."})
		return
	}
	ctx.JSON(200, gin.H{"message": "Usuário atualizado com sucesso."})
}
