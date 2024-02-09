package repository

import "github.com/n9mi/go-course-app/internal/entity"

type PurchaseRepository struct {
	Repository[entity.Purchase]
}

func NewPurchaseRepository() *PurchaseRepository {
	return new(PurchaseRepository)
}
