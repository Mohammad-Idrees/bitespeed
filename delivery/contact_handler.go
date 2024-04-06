package delivery

import (
	"bitespeed/models"
	"bitespeed/repository"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ContactHandler struct {
	repo *repository.ContactRepo
}

const (
	PRIMARY   = "primary"
	SECONDARY = "secondary"
)

var (
	TRUE bool = true
)

func (h *ContactHandler) IdentifyContact(c echo.Context) error {
	ctx := c.Request().Context()
	req := &models.IdentifyContactReq{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// get all contacts with matching email/phone
	getContactParams := &models.GetContactParams{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
	matchingContacts, err := h.repo.GetContactsByEmailOrPhone(ctx, getContactParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// get primaryContactId from matching email/phone
	primaryContactId, err := h.getPrimartContactIdFromMatchingContacts(matchingContacts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// no primary contacts found, insert primary contact and return
	if primaryContactId == 0 {
		contact, err := h.insertPrimaryContact(ctx, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		res := buildIdentifyContactResp(&[]models.Contact{*contact})
		return c.JSON(http.StatusOK, res)
	}

	// get contacts which are primary and should be turned to secondary
	contactIdsToTurnToSecondaryContact := getPrimaryContactsToTurnToSecondary(matchingContacts, primaryContactId)

	// update secondary contacts
	err = h.repo.UpdateLinkedIdAndLinkPrecedenceById(ctx, contactIdsToTurnToSecondaryContact, primaryContactId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// get contacts
	contacts, err := h.repo.GetContactsByPrimaryContactId(ctx, primaryContactId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// matching contact exists with same email and phone from request, return
	if !isNewContactData(req, contacts) {
		res := buildIdentifyContactResp(contacts)
		return c.JSON(http.StatusOK, res)
	}

	// matching contact not found, insert req as secondary contact and return
	linkedId := (*contacts)[0].Id

	// matching contact is secondary
	if (*contacts)[0].LinkedId != nil {
		linkedId = (*contacts)[0].Id
	}

	contact, err := h.insertSecondaryContact(ctx, req, linkedId)
	if err != nil {
		fmt.Println("inside insertSecondaryContact")
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	*contacts = append(*contacts, *contact)
	resp := buildIdentifyContactResp(contacts)

	return c.JSON(http.StatusOK, resp)
}

func (h *ContactHandler) getPrimartContactIdFromMatchingContacts(matchingContacts *[]models.Contact) (int, error) {
	// no matching contacts found
	if len(*matchingContacts) == 0 {
		return 0, nil
	}

	// the top matching contact is secondary? return linkedId
	if (*matchingContacts)[0].LinkedId == nil {
		return (*matchingContacts)[0].Id, nil
	}
	return *(*matchingContacts)[0].LinkedId, nil
}

func getPrimaryContactsToTurnToSecondary(matchingContacts *[]models.Contact, primaryContactId int) *[]int {
	contactIdsToTurnToSecondary := make([]int, 0, len(*matchingContacts))
	for _, contact := range *matchingContacts {
		if contact.LinkedId == nil && contact.Id != primaryContactId {
			contactIdsToTurnToSecondary = append(contactIdsToTurnToSecondary, contact.Id)
		}
	}

	return &contactIdsToTurnToSecondary
}

func (h *ContactHandler) insertPrimaryContact(ctx context.Context, req *models.IdentifyContactReq) (*models.Contact, error) {
	insertContactParams := &models.InsertContactParams{
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		LinkedId:       nil,
		LinkPrecedence: PRIMARY,
	}
	id, err := h.repo.InsertContact(ctx, insertContactParams)
	if err != nil {
		return nil, err
	}
	contact, err := h.repo.GetContactById(ctx, id)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (h *ContactHandler) insertSecondaryContact(ctx context.Context, req *models.IdentifyContactReq, linkedId int) (*models.Contact, error) {
	insertContactParams := &models.InsertContactParams{
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		LinkedId:       &linkedId,
		LinkPrecedence: SECONDARY,
	}
	id, err := h.repo.InsertContact(ctx, insertContactParams)
	if err != nil {
		return nil, err
	}
	contact, err := h.repo.GetContactById(ctx, id)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func isNewContactData(req *models.IdentifyContactReq, contacts *[]models.Contact) bool {
	existingEmails := make(map[string]struct{})
	existingPhones := make(map[string]struct{})

	for _, contact := range *contacts {
		if contact.Email != nil {
			existingEmails[*contact.Email] = struct{}{}
		}
		if contact.PhoneNumber != nil {
			existingPhones[*contact.PhoneNumber] = struct{}{}
		}
	}

	if req.Email != nil {
		if _, present := existingEmails[*req.Email]; !present {
			return true
		}
	}

	if req.PhoneNumber != nil {
		if _, present := existingPhones[*req.PhoneNumber]; !present {
			return true
		}
	}

	return false
}

func buildIdentifyContactResp(contacts *[]models.Contact) *models.IdentifyContactResp {
	// contacts param validations, should have atleast one valid contact
	if contacts == nil || len(*contacts) == 0 {
		return nil
	}

	var primaryContactId int
	uniqueEmails := make(map[string]struct{})
	uniquePhoneNumbers := make(map[string]struct{})
	secondaryContactIds := make([]int, 0, len(*contacts)-1)

	for contactIndex, contact := range *contacts {
		if contact.Email != nil {
			if _, present := uniqueEmails[*contact.Email]; !present {
				uniqueEmails[*contact.Email] = struct{}{}
			}
		}
		if contact.PhoneNumber != nil {
			if _, present := uniquePhoneNumbers[*contact.PhoneNumber]; !present {
				uniquePhoneNumbers[*contact.PhoneNumber] = struct{}{}
			}
		}

		if contactIndex == 0 {
			primaryContactId = contact.Id
		} else {
			secondaryContactIds = append(secondaryContactIds, contact.Id)
		}
	}

	emails := make([]string, 0, len(uniqueEmails))
	for email := range uniqueEmails {
		emails = append(emails, email)
	}

	phoneNumbers := make([]string, 0, len(uniquePhoneNumbers))
	for phoneNumer := range uniquePhoneNumbers {
		phoneNumbers = append(phoneNumbers, phoneNumer)
	}

	identifyContactResp := &models.IdentifyContactResp{
		Contact: models.IdentifyContact{
			PrimaryContactId:    primaryContactId,
			Emails:              emails,
			PhoneNumbers:        phoneNumbers,
			SecondaryContactIds: secondaryContactIds,
		},
	}

	return identifyContactResp
}

func NewContactHandler(repo *repository.ContactRepo) *ContactHandler {
	return &ContactHandler{
		repo,
	}
}
