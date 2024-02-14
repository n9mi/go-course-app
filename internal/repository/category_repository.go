package repository

import (
	"github.com/n9mi/go-course-app/internal/entity"
	"github.com/n9mi/go-course-app/internal/helper"
	"github.com/n9mi/go-course-app/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	Repository[entity.Category]
}

func NewCategoryRepository() *CategoryRepository {
	return new(CategoryRepository)
}

func (r *CategoryRepository) List(tx *gorm.DB, listRequest *model.CategoryListRequest) ([]model.CategoryResponse, error) {
	var categories []model.CategoryResponse

	if listRequest.Page > 0 && listRequest.PageSize > 0 {
		tx = tx.Scopes(helper.Paginate(listRequest.Page, listRequest.PageSize))
	}
	query := tx.Model(new(entity.Category)).
		Select(`categories.id,
			categories.name,
			categories.created_at,
			categories.updated_at,
			users.name as created_by,
			categories.member_count`).
		Joins("inner join users on users.id = categories.created_by")

	if len(listRequest.UserID) > 0 {
		query = query.Where("categories.created_by = ?", listRequest.UserID)
	}

	if listRequest.SortByPopular {
		query = query.Order("member_count DESC")
	} else {
		query = query.Order("created_at ASC")
	}

	if err := query.Scan(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) FindByID(tx *gorm.DB, category *model.CategoryResponse, ID string) error {
	query := tx.Model(new(entity.Category)).
		Select(`categories.id,
			categories.name,
			categories.created_at,
			categories.updated_at,
			users.name as created_by,
			categories.member_count`).
		Joins("inner join users on users.id = categories.created_by").
		Where("categories.id = ?", ID)

	if err := query.Scan(&category).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) FindByIDandUserID(tx *gorm.DB, category *entity.Category, ID string, userID string) error {
	return tx.Where("id = ? and created_by = ?", ID, userID).Take(category).Error
}
