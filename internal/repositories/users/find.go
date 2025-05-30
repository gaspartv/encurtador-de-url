package repositories

import (
	"github.com/gaspartv/encurtador-de-url/db"
	"github.com/gaspartv/encurtador-de-url/internal/models"
)

func FindUserRepository(userEmail string) (models.User, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return models.User{}, err
	}
	defer conn.Close()

	query := "SELECT id, email, password FROM users WHERE email LIKE $1 LIMIT 1"
	var user models.User
	err = conn.QueryRow(query, userEmail).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
