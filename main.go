package main

import (
	"bitespeed/config"
	"bitespeed/delivery"
	"bitespeed/models"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type App struct {
	cfg   *config.StartupConfig
	_echo *echo.Echo
}

func newApp(cfg *config.StartupConfig, _echo *echo.Echo) *App {
	return &App{
		cfg:   cfg,
		_echo: _echo,
	}
}

func newRouter(helloHandler *delivery.HelloHandler) *echo.Echo {
	e := echo.New()
	e.GET("/", helloHandler.HelloWorld)
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
	app._echo.Start(app.cfg.Server.Address)

}

func newDatabase(cfg *config.StartupConfig) (*models.Database, error) {
	dbCfg := cfg.Database
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbCfg.User, dbCfg.Pass, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	dbConn, err := sqlx.Open("mysql", conn)
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
