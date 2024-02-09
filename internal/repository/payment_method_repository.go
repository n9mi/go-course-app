package repository

import "github.com/n9mi/go-course-app/internal/entity"

type PaymentMethodRepository struct {
	Repository[entity.PaymentMethod]
}

func NewPaymentMethodRepository() *PaymentMethodRepository {
	return new(PaymentMethodRepository)
}
