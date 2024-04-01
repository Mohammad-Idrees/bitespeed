package delivery

import (
	"bitespeed/repository"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HelloHandler struct {
	repo *repository.HelloRepo
}

func (h *HelloHandler) HelloWorld(c echo.Context) error {
	ctx := c.Request().Context()
	err := h.repo.Ping(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("error: %s, connecting to db", err.Error()))
	}
	return c.JSON(http.StatusOK, "Hello World")
}

func NewHelloHandler(repo *repository.HelloRepo) *HelloHandler {
	return &HelloHandler{
		repo,
	}
}
