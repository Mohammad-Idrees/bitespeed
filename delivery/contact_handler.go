package delivery

import (
	"bitespeed/models"
	"bitespeed/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ContactHandler struct {
	repo *repository.ContactRepo
}

func (h *ContactHandler) GetContact(c echo.Context) error {
	ctx := c.Request().Context()
	req := models.GetContactReq{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	getContactParams := &models.GetContactParams{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
	contacts, err := h.repo.GetContact(ctx, getContactParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(*contacts) > 0 {
		return c.JSON(http.StatusOK, contacts)
	}

	// no contact exists
	insertContactParams := &models.InsertContactParams{
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		LinkedId:       nil,
		LinkPrecedence: "primary",
	}
	id, err := h.repo.InsertContact(ctx, insertContactParams)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}

func NewContactHandler(repo *repository.ContactRepo) *ContactHandler {
	return &ContactHandler{
		repo,
	}
}
