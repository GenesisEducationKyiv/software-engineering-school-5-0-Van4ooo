package repositories

import (
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *SubscriptionRepository) Confirm(token string) error {
	return r.db.Model(&models.Subscription{}).
		Where("token = ?", token).
		Update("confirmed", true).Error
}

func (r *SubscriptionRepository) Delete(token string) error {
	return r.db.Where("token = ?", token).
		Delete(&models.Subscription{}).Error
}
