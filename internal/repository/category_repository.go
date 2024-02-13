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

func (r *CategoryRepository) List(tx *gorm.DB, request *model.CategoryListRequest) ([]model.CategoryResponse, error) {
	var categories []model.CategoryResponse

	if request.Page > 0 && request.PageSize > 0 {
		tx = tx.Scopes(helper.Paginate(request.Page, request.PageSize))
	}
	query := tx.Model(new(entity.Category)).
		Select(`categories.id,
			categories.name,
			categories.created_at,
			categories.updated_at,
			users.name as created_by,
			categories.member_count`).
		Joins("inner join users on users.id = categories.created_by")

	if len(request.UserID) > 0 {
		query = query.Where("categories.created_by = ?", request.UserID)
	}

	if request.SortByPopular {
		query = query.Order("member_count DESC")
	} else {
		query = query.Order("created_at ASC")
	}

	if err := query.Scan(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
