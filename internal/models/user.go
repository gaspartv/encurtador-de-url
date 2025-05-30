package models

import "time"

type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Disabled  bool       `json:"disabled"`
}
