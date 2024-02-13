package entity

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	ID            string `gorm:"primaryKey"`
	Name          string
	Description   string
	CategoryID    string
	PriceIdr      float64
	BannerLink    string
	CreatedBy     string
	MemberCount   uint64
	Purchases     []Purchase
	CourseMembers []CourseMember
}
