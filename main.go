package main

import (
	"context"
	"log"
	"task/config"
	"task/database"
)

const configPath = "./config/config.yaml"

func main() {
	cfg := config.MustLoad(configPath)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	dbControl, err := database.NewSQliteFloodControl(cfg.StoragePath, ctx)
	if err != nil {
		log.Fatalf("failed to create new sqlite flood control: %v", err)

	}
	passedCheck, err := dbControl.Check(ctx, 1)
	if err != nil {
		log.Fatalf("failed to check flood control: %v", err)
	}
	if !passedCheck {
		log.Fatalf("flood control check failed")
	}
	if passedCheck {
		log.Println("flood control check passed")
	}
}
