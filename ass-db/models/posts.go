package models

import (
	"time"

	"gorm.io/gorm"
)

type ReviewStatus string

const (
	Pending  ReviewStatus = "Pending"
	Approved ReviewStatus = "Approved"
	Declined ReviewStatus = "Declined"
)

type Posts struct {
	ID           uint      `gorm: "primary key;autoincrement" json:id"`
	Description  *string   `json:"description"`
	ReviewStatus *string   `json:"review_status"`
	UserID       uint      `json:"user_id"`
	AdminID      *uint     `json:"admin_id"`
	ReviewDate   time.Time `json:"review_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         Users     `gorm:"forignKey:UserID"`
	Admin        *Admins   `gorm:"foreignKey:AdminID"`
}

func MigratePosts(db *gorm.DB) error {
	err := db.AutoMigrate(&Posts{})
	return err
}
