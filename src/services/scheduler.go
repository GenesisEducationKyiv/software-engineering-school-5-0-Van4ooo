package services

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type Scheduler struct {
	DB             *gorm.DB
	WeatherService WeatherService
	EmailSender    EmailSender
}

type EmailSender interface {
	Send(to, subject, body string) error
}

func NewScheduler(cfg *config.AppConfig) *Scheduler {
	database := db.DB

	weatherSvc := NewOpenWeatherService(
		cfg.WeatherAPI.Key,
		cfg.WeatherAPI.BaseURL,
	)

	emailSender := NewSMTPEmailSender(
		cfg.SMTP.Host,
		cfg.SMTP.Addr,
		cfg.SMTP.Name,
		cfg.SMTP.Pass,
	)

	return &Scheduler{
		DB:             database,
		WeatherService: weatherSvc,
		EmailSender:    emailSender,
	}
}

func (s *Scheduler) Start() {
	c := cron.New(cron.WithLocation(time.UTC))

	if _, err := c.AddFunc("0 * * * *", func() {
		s.sendWeatherUpdates("hourly")
	}); err != nil {
		log.Fatalf("Error scheduling hourly updates: %v", err)
	}

	if _, err := c.AddFunc("0 8 * * *", func() {
		s.sendWeatherUpdates("daily")
	}); err != nil {
		log.Fatalf("Error scheduling daily updates: %v", err)
	}

	c.Start()
	select {}
}

func (s *Scheduler) sendWeatherUpdates(freq string) {
	var subs []models.Subscription
	r := s.DB.Where("frequency = ? AND confirmed = true", freq).Find(&subs)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		log.Printf("Error fetching subscriptions: %v", r.Error)
		return
	}

	for _, sub := range subs {
		weather, err := s.WeatherService.GetWeather(sub.City)
		if err != nil {
			log.Printf("Error fetching weather for %s: %v", sub.City, err)
			continue
		}

		body := fmt.Sprintf(
			"Current weather in %s: %.1f°C, %s (Humidity %.0f%%)",
			sub.City,
			weather.Temperature,
			weather.Description,
			weather.Humidity,
		)

		if err := s.EmailSender.Send(sub.Email, "Weather Update", body); err != nil {
			log.Printf("Error sending email to %s: %v", sub.Email, err)
		}
	}
}
