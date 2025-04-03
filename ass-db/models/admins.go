package models

import "gorm.io/gorm"

type Admins struct {
	ID        uint    `gorm: "primary key;autoincrement" json:id"`
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	Role      *string `json:"role"`
}

func MigrateAdmins(db *gorm.DB) error {
	err := db.AutoMigrate(&Admins{})
	return err
}
