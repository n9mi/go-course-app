package repository

import (
	"fmt"
	"strings"

	"github.com/n9mi/go-course-app/internal/entity"
	"github.com/n9mi/go-course-app/internal/helper"
	"github.com/n9mi/go-course-app/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return new(UserRepository)
}

func (r *UserRepository) List(tx *gorm.DB, listRequest *model.UserListRequest) ([]entity.User, error) {
	var users []entity.User

	var userIDRoleMatch []string
	if len(listRequest.FilterRoleID) > 0 {
		strQuery := "where exists "
		for i, rID := range listRequest.FilterRoleID {
			strQuery += fmt.Sprintf("(select 1 from user_roles where user_roles.user_id = users.id and user_roles.role_id = '%s')", rID)
			if i < len(listRequest.FilterRoleID)-1 {
				strQuery += " and exists "
			}
		}

		err := tx.Raw("select users.id from users " + strQuery).Scan(&userIDRoleMatch).Error

		if err != nil {
			return nil, err
		}
	}

	if listRequest.Page > 0 && listRequest.PageSize > 0 {
		tx = tx.Scopes(helper.Paginate(listRequest.Page, listRequest.PageSize))
	}

	if len(listRequest.SearchName) > 0 {
		tx = tx.Where("lower(name) like ?", "%"+listRequest.SearchName+"%")
	}

	if len(listRequest.SearchEmail) > 0 {
		tx = tx.Where("lower(email) like ?", "%"+listRequest.SearchEmail+"%")
	}

	var err error
	if len(listRequest.FilterRoleID) > 0 {
		err = tx.Preload("Roles").Find(&users, "id in (?)", userIDRoleMatch).Error
	} else {
		err = tx.Preload("Roles").Find(&users).Error
	}

	return users, err
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
