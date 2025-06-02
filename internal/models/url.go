package models

import "time"

type Url struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	UserID      int64     `json:"user_id" gorm:"not null"`
	OriginalUrl string    `json:"original_url" gorm:"not null"`
	ShortUrl    string    `json:"short_url" gorm:"not null;unique"`
	Clicks      int64     `json:"clicks" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
}
