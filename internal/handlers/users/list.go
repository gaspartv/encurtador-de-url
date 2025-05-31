package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gaspartv/encurtador-de-url/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListUsersHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	userId := ctx.Query("id")
	email := ctx.Query("email")
	name := ctx.Query("name")
	disabled := ctx.Query("disabled")
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("page_number", "1"))
	sortBy := ctx.DefaultQuery("sort_by", "id")
	sortOrder := strings.ToLower(ctx.DefaultQuery("sort_order", "asc"))

	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize

	allowedSortFields := map[string]bool{
		"id":         true,
		"email":      true,
		"name":       true,
		"created_at": true,
		"updated_at": true,
		"disabled":   true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	db := ctx.MustGet("db").(*gorm.DB)

	var users []models.User
	query := db.Model(&models.User{})

	if userId != "" {
		query = query.Where("id = ?", userId)
	}

	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if disabled != "" {
		query = query.Where("disabled = ?", disabled)
	}

	var total int64
	query.Count(&total)

	query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(pageSize).
		Offset(offset).
		Find(&users)

	var totalPages int
	if total > 0 {
		totalPages = (int(total) + pageSize - 1) / pageSize
	} else {
		totalPages = 0
	}

	var nextPage int
	if pageNumber < totalPages {
		nextPage = pageNumber + 1
	} else {
		nextPage = 0
	}

	var previousPage int
	if pageNumber > 1 {
		previousPage = pageNumber - 1
	} else {
		previousPage = 0
	}

	ctx.JSON(http.StatusOK, gin.H{
		"page_size":     pageSize,
		"page_number":   pageNumber,
		"total":         total,
		"total_pages":   totalPages,
		"next_page":     nextPage,
		"previous_page": previousPage,
		"sort_by":       sortBy,
		"sort_order":    sortOrder,
		"data":          users,
	})
}
