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
	if err := tx.Create(&entity.CourseMember{
		CourseID: course.ID,
		UserID:   user.ID,
	}).Error; err != nil {
		return err
	}

	if err := tx.Model(new(entity.Category)).Where("id = ?", course.CategoryID).
		UpdateColumn("member_count", gorm.Expr("member_count + ?", 1)).Error; err != nil {
		return err
	}

	if err := tx.Model(new(entity.Course)).Where("id = ?", course.ID).
		UpdateColumn("member_count", gorm.Expr("member_count + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}
