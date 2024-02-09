package entity

import (
	"time"

	"gorm.io/gorm"
)

type Purchase struct {
	gorm.Model
	ID               string `gorm:"primaryKey"`
	UserID           string
	CourseID         string
	PaymentMethodID  uint64
	Status           uint8
	PurchaseDeadline time.Time
	PurchasedAt      *time.Time
}
