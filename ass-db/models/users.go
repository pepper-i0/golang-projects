package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint      `gorm: "primary key;autoincrement" json:id"`
	Firstname *string   `json:"firstname"`
	Lastname  *string   `json:"lastname"`
	Email     *string   `json:"email"`
	Phone     *int      `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updatedat"`
	// Posts     []Posts   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
