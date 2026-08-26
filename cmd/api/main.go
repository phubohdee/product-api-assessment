package main

import (
	"fmt"
	"log"
	"os"

	"product-api-assessment/internal/config"
	handler "product-api-assessment/internal/handler/http"
	"product-api-assessment/internal/repository/postgres"
	"product-api-assessment/internal/service"
)

// @title          Product API
// @version        1.0
// @description    RESTful API for managing products
// @host           localhost:8080
// @BasePath       /
func main() {
	cfg := config.Load()

	args := os.Args[1:]

	// Handle migrate subcommand: go run cmd/api/main.go migrate up|down
	if len(args) >= 2 && args[0] == "migrate" {
		runMigrate(cfg, args[1])
		return
	}

	// Start API server
	db, err := cfg.DB.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	productRepo := postgres.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r := handler.NewRouter(productHandler)

	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.App.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
