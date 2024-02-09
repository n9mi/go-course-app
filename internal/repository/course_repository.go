package repository

import (
	"github.com/n9mi/go-course-app/internal/entity"
	"gorm.io/gorm"
)

type CourseRepository struct {
	Repository[entity.Course]
}

func NewCourseRepository() *CourseRepository {
	return new(CourseRepository)
}

func (r *CourseRepository) AddMember(tx *gorm.DB, course *entity.Course, user *entity.User) error {
	return tx.Create(&entity.CourseMember{
		CourseID: course.ID,
		UserID:   user.ID,
	}).Error
}
