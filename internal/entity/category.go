package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Name      string
	CreatedBy string
	Courses   []Course
}
