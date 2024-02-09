package repository

import "github.com/n9mi/go-course-app/internal/entity"

type CategoryRepository struct {
	Repository[entity.Category]
}

func NewCategoryRepository() *CategoryRepository {
	return new(CategoryRepository)
}
