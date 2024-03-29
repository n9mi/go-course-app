package repository

import (
	"github.com/n9mi/go-course-app/internal/entity"
	"github.com/n9mi/go-course-app/internal/helper"
	"github.com/n9mi/go-course-app/internal/model"
	"gorm.io/gorm"
)

type CourseRepository struct {
	Repository[entity.Course]
}

func NewCourseRepository() *CourseRepository {
	return new(CourseRepository)
}

func (r *CourseRepository) List(tx *gorm.DB, listRequest *model.CourseListRequest) ([]model.CourseListResponse, error) {
	var courses []model.CourseListResponse

	if listRequest.Page > 0 && listRequest.PageSize > 0 {
		tx = tx.Scopes(helper.Paginate(listRequest.Page, listRequest.PageSize))
	}

	tx = tx.Model(new(entity.Course)).
		Select(`courses.id,
			courses.name,
			categories.name as category,
			courses.price_idr,
			courses.banner_link,
			users.name as created_by,
			courses.member_count`).
		Joins("inner join categories on categories.id = courses.category_id").
		Joins("inner join users on users.id = courses.created_by")

	if len(listRequest.UserID) > 0 {
		tx = tx.Where("courses.created_by = ?", listRequest.UserID)
	}

	if len(listRequest.CategoryID) > 0 {
		tx = tx.Where("courses.category_id = ?", listRequest.CategoryID)
	}

	if len(listRequest.SearchTitle) > 0 {
		tx = tx.Where("lower(courses.name) like ?", "%"+listRequest.SearchTitle+"%")
	}

	if listRequest.IsFree {
		tx = tx.Where("courses.price_idr = ?", 0)
	} else if listRequest.SortByMaximumPrice {
		tx = tx.Order("courses.price_idr desc")
	} else if listRequest.SortByMinimumPrice {
		tx = tx.Order("courses.price_idr asc")
	}

	if err := tx.Scan(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepository) ScanByIDAndUserID(tx *gorm.DB, course *model.CourseResponse, ID string, userID string) error {
	tx = tx.Model(new(entity.Course)).
		Select(`courses.id,
			courses.name,
			courses.description,
			categories.name as category,
			courses.price_idr,
			courses.banner_link,
			users.name as created_by,
			courses.member_count`).
		Joins("inner join categories on categories.id = courses.category_id").
		Joins("inner join users on users.id = courses.created_by")

	if len(userID) > 0 {
		tx = tx.Where("courses.id = ? and courses.created_by = ?", ID, userID)
	} else {
		tx = tx.Where("courses.id = ?", ID)
	}
	tx = tx.Scan(course)

	return tx.Error
}

func (r *CourseRepository) FindByIDAndUserID(tx *gorm.DB, course *entity.Course, ID string, userID string) error {
	return tx.Where("id = ? and created_by = ?", ID, userID).Take(course).Error
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
