package model

import "time"

type User struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	Email string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	IsDeleted bool `gorm:"default:false" json:"is_deleted"`
	
	Addresses []Address `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}