// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"bitespeed/config"
	"bitespeed/delivery"
	"bitespeed/repository"
)

// Injectors from wire.go:

func InitializeDependency(cfg *config.StartupConfig) (*App, error) {
	database, err := newDatabase(cfg)
	if err != nil {
		return nil, err
	}
	helloRepo := repository.NewHelloRepo(database)
	helloHandler := delivery.NewHelloHandler(helloRepo)
	echo := newRouter(helloHandler)
	app := newApp(cfg, echo)
	return app, nil
}