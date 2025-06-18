package repositories

import (
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type EmailerSubscriptionStore struct {
	DB *gorm.DB
}

func NewSubscriptionStore(db *gorm.DB) SubscriptionStore {
	return &EmailerSubscriptionStore{DB: db}
}

func (r *EmailerSubscriptionStore) FetchByFrequency(freq string) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := r.DB.Where("frequency = ? AND confirmed = true", freq).Find(&subs).Error
	return subs, err
}
