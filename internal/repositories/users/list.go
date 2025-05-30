package repositories

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gaspartv/encurtador-de-url/db"
	"github.com/gaspartv/encurtador-de-url/internal/models"
)

func ListUsersRepository(
	userID string,
	email string,
	name string,
	disabled string,
	pageSize string,
	pageNumber string,
	SortBy string,
	SortOrder string,
) (models.Pagination, error) {
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil && pageSize != "" {
		return models.Pagination{}, fmt.Errorf("invalid page_size: %v", err)
	}

	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil && pageNumber != "" {
		return models.Pagination{}, fmt.Errorf("invalid page_number: %v", err)
	}

	if pageSizeInt <= 0 {
		pageSizeInt = 10
	}

	if pageNumberInt <= 0 {
		pageNumberInt = 1
	}

	conn, err := db.OpenConnection()
	if err != nil {
		return models.Pagination{}, err
	}
	defer conn.Close()

	query := `SELECT id, email, password, name, created_at, updated_at, deleted_at, disabled FROM users`
	queryCount := "SELECT COUNT(*) FROM users"
	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if userID != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argID))
		args = append(args, userID)
		argID++
	}

	if email != "" {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argID))
		args = append(args, "%"+email+"%")
		argID++
	}

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argID))
		args = append(args, "%"+name+"%")
		argID++
	}

	if disabled != "" {
		conditions = append(conditions, fmt.Sprintf("disabled = $%d", argID))
		args = append(args, disabled)
		argID++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
		queryCount += " WHERE " + strings.Join(conditions, " AND ")
	}

	validSortBy := map[string]bool{"id": true, "email": true, "name": true, "created_at": true, "updated_at": true, "disabled": true}
	validSortOrder := map[string]bool{"asc": true, "desc": true}

	sortBy := "id"
	sortOrder := "asc"

	if SortBy != "" && validSortBy[SortBy] {
		sortBy = SortBy
	}
	if SortOrder != "" && validSortOrder[strings.ToLower(SortOrder)] {
		sortOrder = strings.ToLower(SortOrder)
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	query += fmt.Sprintf(" LIMIT %d", pageSizeInt)

	if pageNumberInt > 1 {
		offset := (pageNumberInt - 1) * pageSizeInt
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	rows, err := conn.Query(query, args...)
	if err != nil {
		return models.Pagination{}, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Disabled,
		)
		if err != nil {
			return models.Pagination{}, err
		}

		users = append(users, user)
	}

	var totalCount int
	err = conn.QueryRow(queryCount, args...).Scan(&totalCount)
	if err != nil {
		return models.Pagination{}, err
	}

	var totalPages int
	if totalCount > 0 {
		totalPages = (totalCount + pageSizeInt - 1) / pageSizeInt
	} else {
		totalPages = 0
	}

	var nextPage int
	if pageNumberInt < totalPages {
		nextPage = pageNumberInt + 1
	} else {
		nextPage = 0
	}

	var previousPage int
	if pageNumberInt > 1 {
		previousPage = pageNumberInt - 1
	} else {
		previousPage = 0
	}

	return models.Pagination{
		PageSize:     pageSizeInt,
		PageNumber:   pageNumberInt,
		TotalCount:   totalCount,
		TotalPages:   totalPages,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Data:         users,
	}, nil
}
