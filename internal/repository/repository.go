package repository

import "gorm.io/gorm"

type Repository[T any] struct {
}

func (r *Repository[T]) FindByID(tx *gorm.DB, e *T, ID any) error {
	return tx.First(e, "id = ?", ID).Error
}

func (r *Repository[T]) Create(tx *gorm.DB, e *T) error {
	return tx.Create(e).Error
}

func (r *Repository[T]) Updates(tx *gorm.DB, e *T) error {
	return tx.Model(e).Updates(e).Error
}

func (r *Repository[T]) Save(tx *gorm.DB, e *T) error {
	return tx.Save(e).Error
}

func (r *Repository[T]) CountByID(tx *gorm.DB, ID any) (int64, error) {
	var count int64
	err := tx.Model(new(T)).Where("id = ?", ID).Count(&count).Error

	return count, err
}

func (r *Repository[T]) Delete(tx *gorm.DB, e *T) error {
	return tx.Delete(e).Error
}
