package entity

import "gorm.io/gorm"

type CourseMember struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey"`
	UserID   string
	CourseID string
}
