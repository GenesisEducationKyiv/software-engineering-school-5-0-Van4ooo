package repositories

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"

	"gorm.io/gorm"
)

type GormSubscription struct {
	db *gorm.DB
}

func NewGormSubscription(db *gorm.DB) *GormSubscription {
	return &GormSubscription{db: db}
}

func (r *GormSubscription) Save(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *GormSubscription) MarkConfirmed(token string) error {
	return r.db.Model(&models.Subscription{}).
		Where("token = ?", token).
		Update("confirmed", true).
		Error
}

func (r *GormSubscription) Remove(token string) error {
	return r.db.Where("token = ?", token).
		Delete(&models.Subscription{}).
		Error
}
