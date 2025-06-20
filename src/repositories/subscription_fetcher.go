package repositories

import (
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type SubcsriptionFetcher struct {
	DB *gorm.DB
}

func NewSubscriptionStore(db *gorm.DB) *SubcsriptionFetcher {
	return &SubcsriptionFetcher{DB: db}
}

func (r *SubcsriptionFetcher) FetchByFrequency(
	freq string) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := r.DB.Where("frequency = ? AND confirmed = true", freq).Find(&subs).Error
	return subs, err
}
