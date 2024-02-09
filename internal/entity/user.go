package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             string `gorm:"primaryKey"`
	Name           string
	Email          string
	Password       string
	Roles          []*Role    `gorm:"many2many:user_roles"`
	Categories     []Category `gorm:"foreignKey:CreatedBy"`
	Purchases      []Purchase
	Courses        []Course `gorm:"foreignKey:CreatedBy"`
	CourseMembers  []CourseMember
	PaymentMethods []PaymentMethod `gorm:"foreignKey:CreatedBy"`
}
