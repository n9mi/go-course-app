package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Name        string
	CreatedBy   string
	MemberCount uint64
	Courses     []Course
}
