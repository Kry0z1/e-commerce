package main

import (
	"fmt"

	"github.com/Kry0z1/e-commerce/sso-microservice/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: Configure logger
	// TODO: Start connection

	// TODO: Initialize service

	// TODO: Graceful shutdown
}
