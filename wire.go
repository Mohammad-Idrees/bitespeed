//go:build wireinject
// +build wireinject

package main

import (
	"bitespeed/config"
	delivery "bitespeed/delivery"
	repository "bitespeed/repository"

	"github.com/google/wire"
)

func InitializeDependency(cfg *config.StartupConfig) (*App, error) {
	wire.Build(
		newDatabase,
		repository.NewContactRepo,
		delivery.NewContactHandler,
		newRouter,
		newApp,
	)
	return &App{}, nil
}
