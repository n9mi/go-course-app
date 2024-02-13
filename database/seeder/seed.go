package seeder

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/n9mi/go-course-app/internal/entity"
	"github.com/n9mi/go-course-app/internal/helper"
	"github.com/n9mi/go-course-app/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, redisClient *redis.Client, repositoryBootstrap *repository.RepositorySetup) error {
	roles, err := SeedRoles(db, repositoryBootstrap.RoleRepository)
	if err != nil {
		return err
	}

	fmt.Println("Successfully seeds the roles")

	users, err := SeedUsers(db, repositoryBootstrap.UserRepository, roles)
	if err != nil {
		return err
	}

	fmt.Println("Successfully seeds the users")

	categories, err := SeedCategories(db, repositoryBootstrap.CategoryRepository,
		repositoryBootstrap.UserRepository, users)
	if err != nil {
		return err
	}

	fmt.Println("Successfully seeds the categories")

	courses, err := SeedCourses(db, repositoryBootstrap.CourseRepository, categories)
	if err != nil {
		return err
	}

	fmt.Println("Successfully seeds the courses")

	paymentMethods, err := SeedPaymentMethod(db, repositoryBootstrap.PaymentMethodRepository,
		repositoryBootstrap.UserRepository, users)
	if err != nil {
		return err
	}

	fmt.Println("Successfully seeds the payment methods")

	purchases, err := SeedPurchases(db, repositoryBootstrap.PurchaseRepository, repositoryBootstrap.UserRepository,
		users, courses, paymentMethods)
	if err != nil {
		return err
	}

	fmt.Println("Successfully seeds the purchases")

	err = SeedCourseMembers(db, repositoryBootstrap.CourseRepository, repositoryBootstrap.UserRepository,
		purchases)
	if err != nil {
		return err
	}

	fmt.Println("Successfully seeds the course members")

	return nil
}

func SeedRoles(db *gorm.DB, roleRepository *repository.RoleRepository) ([]entity.Role, error) {
	roles := []entity.Role{
		{ID: "admin", DisplayName: "Admin"},
		{ID: "user", DisplayName: "User"},
	}

	for _, r := range roles {
		tx := db.Begin()
		if err := roleRepository.Repository.Save(tx, &r); err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
	}

	return roles, nil
}

func SeedUsers(db *gorm.DB, userRepository *repository.UserRepository, roles []entity.Role) ([]entity.User, error) {
	var users []entity.User

	for _, r := range roles {
		for i := 1; i <= 3; i++ {
			newPassword, _ := helper.GeneratePassword("password")
			newUser := entity.User{
				ID:       "USR_" + helper.GenerateRandomString(12),
				Name:     fmt.Sprintf("%s %d", r.DisplayName, i),
				Email:    fmt.Sprintf("%s%d@mail.com", r.ID, i),
				Password: newPassword,
			}
			tx := db.Begin()
			if err := userRepository.Repository.Save(tx, &newUser); err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := userRepository.AssignRoles(tx, &newUser, []entity.Role{r}); err != nil {
				tx.Rollback()
				return nil, err
			}
			tx.Commit()
			users = append(users, newUser)
		}
	}

	return users, nil
}

func SeedCategories(db *gorm.DB, categoryRepository *repository.CategoryRepository,
	userRepository *repository.UserRepository, users []entity.User) ([]entity.Category, error) {
	var categories []entity.Category

	for _, u := range users {
		tx := db.Begin()
		if userRepository.HasRole(tx, &u, "admin") {
			for i := 1; i <= 3; i++ {
				tx = db.Begin()
				newCategory := entity.Category{
					ID:        "CAT_" + helper.GenerateRandomString(10),
					Name:      fmt.Sprintf("Category %s", helper.GenerateRandomString(4)),
					CreatedBy: u.ID,
				}
				if err := categoryRepository.Save(tx, &newCategory); err != nil {
					tx.Rollback()
					return nil, err
				}
				tx.Commit()
				categories = append(categories, newCategory)
			}
		}
	}

	return categories, nil
}

func SeedCourses(db *gorm.DB, courseRepository *repository.CourseRepository,
	categories []entity.Category) ([]entity.Course, error) {
	var courses []entity.Course
	minPrice := 10000
	maxPrice := 100000

	for _, c := range categories {
		for i := 0; i < 2; i++ {
			tx := db.Begin()
			newCourse := entity.Course{
				ID:          helper.GenerateRandomString(16),
				Name:        "Course " + helper.GenerateRandomString(25),
				Description: helper.GenerateRandomString(50),
				CategoryID:  c.ID,
				PriceIdr:    float64(rand.Intn(maxPrice-minPrice) + minPrice),
				BannerLink:  "https://picsum.photos/400",
				CreatedBy:   c.CreatedBy,
			}
			if err := courseRepository.Repository.Save(tx, &newCourse); err != nil {
				tx.Rollback()
				return nil, err
			}
			tx.Commit()
			courses = append(courses, newCourse)
		}
	}

	return courses, nil
}

func SeedPaymentMethod(db *gorm.DB, paymentMethodRepository *repository.PaymentMethodRepository, userRepository *repository.UserRepository,
	users []entity.User) ([]entity.PaymentMethod, error) {
	var paymentMethod []entity.PaymentMethod

	for _, u := range users {
		tx := db.Begin()
		if userRepository.HasRole(tx, &u, "admin") {
			for i := 1; i <= 2; i++ {
				tx = db.Begin()
				newPaymentMethod := entity.PaymentMethod{
					Name:                 fmt.Sprintf("Payment Method %s", helper.GenerateRandomString(4)),
					Description:          helper.GenerateRandomString(10),
					IsCurrentlyAvailable: true,
					CreatedBy:            u.ID,
				}
				if err := paymentMethodRepository.Repository.Save(tx, &newPaymentMethod); err != nil {
					tx.Rollback()
					return nil, err
				}
				tx.Commit()
				paymentMethod = append(paymentMethod, newPaymentMethod)
			}
		}
	}

	return paymentMethod, nil
}

func SeedPurchases(db *gorm.DB, purchaseRepository *repository.PurchaseRepository, userRepository *repository.UserRepository,
	users []entity.User, courses []entity.Course, paymentMethods []entity.PaymentMethod) ([]entity.Purchase, error) {
	var purchases []entity.Purchase
	for _, u := range users {
		tx := db.Begin()
		if userRepository.HasRole(tx, &u, "user") {
			curCourses := make(map[string]bool)
			for i := 1; i <= 10; i++ {
				tx = db.Begin()
				courseIDSelected := courses[rand.Intn(len(courses))].ID
				_, alreadySelected := curCourses[courseIDSelected] // _, ok
				if !alreadySelected {
					tx = db.Begin()
					purchasedAt := time.Now().Add(time.Duration(rand.Intn(23)) * time.Hour)
					newPurchase := entity.Purchase{
						ID:               helper.GenerateRandomString(10),
						UserID:           u.ID,
						CourseID:         courseIDSelected,
						PaymentMethodID:  paymentMethods[rand.Intn(len(paymentMethods))].ID,
						Status:           3,
						PurchaseDeadline: time.Now().Add(time.Duration(24) * time.Hour),
						PurchasedAt:      &purchasedAt,
					}
					if err := purchaseRepository.Repository.Save(tx, &newPurchase); err != nil {
						tx.Rollback()
						return nil, err
					}
					tx.Commit()
					purchases = append(purchases, newPurchase)
					curCourses[courseIDSelected] = true
				}
			}
		}
	}

	return purchases, nil
}

func SeedCourseMembers(db *gorm.DB, courseRepository *repository.CourseRepository,
	userRepository *repository.UserRepository, purchases []entity.Purchase) error {
	for _, p := range purchases {
		tx := db.Begin()
		courseFound := new(entity.Course)
		courseRepository.Repository.FindByID(tx, courseFound, p.CourseID)
		userFound := new(entity.User)
		userRepository.Repository.FindByID(tx, userFound, p.UserID)

		tx = db.Begin()
		if err := courseRepository.AddMember(tx, courseFound, userFound); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
}
