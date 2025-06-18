package repositories

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"

	"gorm.io/gorm"
)

type GormSubscriptionStorage struct {
	db *gorm.DB
}

func NewGormSubscriptionStorage(db *gorm.DB) SubscriptionStorage {
	return &GormSubscriptionStorage{db: db}
}

func (r *GormSubscriptionStorage) Save(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *GormSubscriptionStorage) MarkConfirmed(token string) error {
	return r.db.Model(&models.Subscription{}).
		Where("token = ?", token).
		Update("confirmed", true).
		Error
}

func (r *GormSubscriptionStorage) Remove(token string) error {
	return r.db.Where("token = ?", token).
		Delete(&models.Subscription{}).
		Error
}
