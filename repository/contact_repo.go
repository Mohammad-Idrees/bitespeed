package repository

import (
	"bitespeed/models"
	"context"
	"fmt"
	"time"
)

type ContactRepo struct {
	db *models.Database
}

const (
	GetContactsByEmailOrPhone                  = "SELECT * FROM Contacts WHERE email = ? OR phoneNumber = ? ORDER BY id ASC"
	InsertContact                              = "INSERT INTO Contacts(phoneNumber, email, linkedId, linkPrecedence) VALUES(?, ?, ?, ?)"
	GetContactById                             = "SELECT * FROM Contacts WHERE id = ?"
	GetContactsByPrimaryContactId              = "SELECT * FROM Contacts WHERE id = ? or linkedId = ? ORDER by id ASC"
	UpdateContactLinkedIdAndLinkPrecedenceById = "UPDATE Contacts SET linkedId = ?, linkPrecedence = ? where id = ?"
	SECONDARY                                  = "secondary"
)

func (r *ContactRepo) UpdateLinkedIdAndLinkPrecedenceById(ctx context.Context, ids *[]int, linkedId int) error {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	// begin transaction
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// prepare query
	stmt, err := tx.PrepareContext(ctx, UpdateContactLinkedIdAndLinkPrecedenceById)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range *ids {
		// execute query
		_, err = stmt.ExecContext(ctx, linkedId, SECONDARY, id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (r *ContactRepo) GetContactsByEmailOrPhone(ctx context.Context, params *models.GetContactParams) (*[]models.Contact, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	contacts := &[]models.Contact{}
	err := r.db.DB.SelectContext(ctx, contacts, GetContactsByEmailOrPhone, params.Email, params.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *ContactRepo) GetContactsByPrimaryContactId(ctx context.Context, primaryContactId int) (*[]models.Contact, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	contacts := &[]models.Contact{}
	err := r.db.DB.SelectContext(ctx, contacts, GetContactsByPrimaryContactId, primaryContactId, primaryContactId)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *ContactRepo) GetContactById(ctx context.Context, id int) (*models.Contact, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	res := []models.Contact{}
	err := r.db.DB.SelectContext(ctx, &res, GetContactById, id)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("contact not found")
	}
	return &res[0], nil
}

func (r *ContactRepo) InsertContact(ctx context.Context, params *models.InsertContactParams) (int, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	// begin transaction
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	// prepare query
	stmt, err := tx.PrepareContext(ctx, InsertContact)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.ExecContext(ctx, params.PhoneNumber, params.Email, params.LinkedId, params.LinkPrecedence)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func NewContactRepo(db *models.Database) *ContactRepo {
	return &ContactRepo{
		db,
	}
}
