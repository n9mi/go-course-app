package entity

import "gorm.io/gorm"

type PaymentMethod struct {
	gorm.Model
	ID                   uint64 `gorm:"primaryKey"`
	Name                 string
	Description          string
	IsCurrentlyAvailable bool
	CreatedBy            string
	Purchases            []Purchase
}
