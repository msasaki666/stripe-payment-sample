package main

import (
	"fmt"

	"github.com/labstack/echo/v4/middleware"
	"github.com/msasaki666/backend/internal/renv"
)

func createCorsConfig(cfg *config) (middleware.CORSConfig, error) {
	corsConfig := middleware.DefaultCORSConfig
	if cfg.GoEnv == renv.Development {
		corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	} else if cfg.GoEnv == renv.Staging {
		// corsConfig.AllowOrigins = []string{"https://staging.example.com"}
	} else if cfg.GoEnv == renv.Production {
		// corsConfig.AllowOrigins = []string{"https://example.com"}
	} else {
		return middleware.CORSConfig{}, fmt.Errorf("invalid env: %s", cfg.GoEnv)
	}
	return corsConfig, nil
}
