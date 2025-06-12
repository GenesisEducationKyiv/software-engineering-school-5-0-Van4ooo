package services

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

func StartScheduler() {
	c := cron.New(cron.WithLocation(time.UTC))

	if _, err := c.AddFunc("0 * * * *", func() {
		sendWeatherUpdates("hourly")
	}); err != nil {
		log.Fatalf("Error adding hourly updates: %v", err)
	}

	if _, err := c.AddFunc("0 8 * * *", func() {
		sendWeatherUpdates("daily")
	}); err != nil {
		log.Fatalf("Error adding daily updates: %v", err)
	}

	c.Start()
	select {}
}

func sendWeatherUpdates(freq string) {
	var subs []models.Subscription
	result := db.DB.Where("frequency = ? AND confirmed = true", freq).Find(&subs)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Error fetching subscriptions: %v", result.Error)
		return
	}
	for _, sub := range subs {
		weather, err := FetchCurrentWeather(sub.City)
		if err != nil {
			log.Printf("Error fetching weather for %s: %v", sub.City, err)
			continue
		}
		body := fmt.Sprintf("Current weather in %s: %.1f°C, %s (Humidity %.0f%%)",
			sub.City, weather.Temperature, weather.Description, weather.Humidity)
		if err := SendEmail(sub.Email, "Weather Update", body); err != nil {
			log.Printf("Error sending email to %s: %v", sub.Email, err)
		}
	}
}
