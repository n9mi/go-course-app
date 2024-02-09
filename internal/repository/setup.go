package repository

type RepositorySetup struct {
	RoleRepository          *RoleRepository
	UserRepository          *UserRepository
	CategoryRepository      *CategoryRepository
	CourseRepository        *CourseRepository
	PaymentMethodRepository *PaymentMethodRepository
	PurchaseRepository      *PurchaseRepository
	CourseMemberRepository  *CourseMemberRepository
}

func Setup() *RepositorySetup {
	return &RepositorySetup{
		RoleRepository:          NewRoleRepository(),
		UserRepository:          NewUserRepository(),
		CategoryRepository:      NewCategoryRepository(),
		CourseRepository:        NewCourseRepository(),
		PaymentMethodRepository: NewPaymentMethodRepository(),
		PurchaseRepository:      NewPurchaseRepository(),
		CourseMemberRepository:  NewCourseMemberRepository(),
	}
}
