package model

import "time"

type Address struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint64 `gorm:"index;not null" json:"user_id"`

	FirstName string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(100)" json:"last_name"`
	Email     string `gorm:"type:varchar(255);index" json:"email"`
	Phone     string `gorm:"type:varchar(20)" json:"phone" validate:"omitempty,phone"`

	AddressLine1 string `gorm:"type:varchar(255)" json:"address_line1"`
	AddressLine2 string `gorm:"type:varchar(255)" json:"address_line2"`
	City         string `gorm:"type:varchar(100)" json:"city"`
	State        string `gorm:"type:varchar(100)" json:"state"`
	Country      string `gorm:"type:varchar(100)" json:"country"`
	Pincode      string `gorm:"type:varchar(20)" json:"pincode"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`
}
