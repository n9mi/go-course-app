package repository

import "github.com/n9mi/go-course-app/internal/entity"

type RoleRepository struct {
	Repository[entity.Role]
}

func NewRoleRepository() *RoleRepository {
	return new(RoleRepository)
}
