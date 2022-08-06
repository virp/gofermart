package config

import (
	"errors"
	"flag"
	"os"
)

type GofermartConfig struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func NewGofermartConfig() (*GofermartConfig, error) {
	cfg := GofermartConfig{}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Run address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "Accrual system address")

	if a, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		cfg.RunAddress = a
	}
	if d, ok := os.LookupEnv("DATABASE_URI"); ok {
		cfg.DatabaseURI = d
	}
	if r, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		cfg.AccrualSystemAddress = r
	}

	if cfg.RunAddress == "" {
		return nil, errors.New("run address is empty")
	}
	if cfg.DatabaseURI == "" {
		return nil, errors.New("database URI is empty")
	}
	if cfg.AccrualSystemAddress == "" {
		return nil, errors.New("accrual system address is empty")
	}

	return &cfg, nil
}
