package repository

import (
	"strings"

	"github.com/n9mi/go-course-app/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return new(UserRepository)
}

func (r *UserRepository) FindByEmail(tx *gorm.DB, user *entity.User, email string) error {
	return tx.First(user, "email = ?", strings.ToLower(email)).Error
}

func (r *UserRepository) HasRole(tx *gorm.DB, user *entity.User, roleID string) bool {
	return tx.Model(user).Where("id = ?", roleID).Association("Roles").Count() > 0
}

func (r *UserRepository) AssignRoles(tx *gorm.DB, user *entity.User, roles []entity.Role) error {
	return tx.Model(user).Association("Roles").Append(roles)
}

func (r *UserRepository) GetRoles(tx *gorm.DB, user *entity.User) ([]entity.Role, error) {
	var rolesFound []entity.Role
	err := tx.Model(user).Association("Roles").Find(&rolesFound)

	return rolesFound, err
}
