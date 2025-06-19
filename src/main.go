package main

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/routers"
)

func main() {
	cfg, err := config.Config()
	if err != nil {
		log.Fatal(err)
	}

	db.Init(cfg.GetDB())

	r := gin.Default()
	routers.SetupRoutes(r, cfg, db.DB)

	services.RunScheduler(cfg, db.DB)

	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
