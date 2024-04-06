package main

import (
	"bitespeed/config"
	"bitespeed/delivery"
	"bitespeed/models"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type App struct {
	cfg  *config.StartupConfig
	echo *echo.Echo
}

func newApp(cfg *config.StartupConfig, echo *echo.Echo) *App {
	return &App{
		cfg:  cfg,
		echo: echo,
	}
}

func newRouter(contactHandler *delivery.ContactHandler) *echo.Echo {
	e := echo.New()
	e.POST("/identify", contactHandler.IdentifyContact)
	return e
}

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed loading config file", err)
		return
	}

	app, err := InitializeDependency(config)
	if err != nil {
		panic(err.Error())
	}

	// start server
	app.echo.Start(app.cfg.Server.Address)

}

func newDatabase(cfg *config.StartupConfig) (*models.Database, error) {
	dbCfg := cfg.Database
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbCfg.User, dbCfg.Pass, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	fmt.Println(url)
	dbConn, err := sqlx.Open("postgres", url)
	if err != nil {
		log.Fatalln("error connecting to database", err.Error())
		return nil, err
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatalln("error pinging databse", err.Error())
		return nil, err
	}
	return &models.Database{DB: dbConn}, nil
}
