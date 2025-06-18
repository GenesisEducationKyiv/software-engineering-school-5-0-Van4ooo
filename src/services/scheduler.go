package services

import (
	"fmt"
	"log"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/repositories"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	store          repositories.SubscriptionStore
	WeatherService WeatherService
	EmailSender    EmailSender
}

type EmailSender interface {
	Send(to, subject, body string) error
	SendConfirmation(to, baseURL, token string) error
}

func NewScheduler(
	store repositories.SubscriptionStore,
	service WeatherService,
	sender EmailSender) *Scheduler {
	return &Scheduler{
		store:          store,
		WeatherService: service,
		EmailSender:    sender,
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
	subs, err := s.store.FetchByFrequency(freq)
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
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
