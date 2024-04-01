package repository

import (
	"bitespeed/models"
	"context"
	"time"
)

type ContactRepo struct {
	db *models.Database
}

const (
	GetContact    = "SELECT * FROM Contacts WHERE email = ? OR phoneNumber = ? ORDER BY id ASC"
	InsertContact = "INSERT INTO Contacts(phoneNumber, email, linkedId, linkPrecedence) VALUES(?, ?, ?, ?)"
)

func (r *ContactRepo) GetContact(ctx context.Context, params *models.GetContactParams) (*[]models.Contact, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	contacts := &[]models.Contact{}
	err := r.db.DB.SelectContext(ctx, contacts, GetContact, params.Email, params.PhoneNumber)
	if err != nil {
		return contacts, err
	}
	return contacts, nil
}

func (r *ContactRepo) InsertContact(ctx context.Context, params *models.InsertContactParams) (int, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	stmt, err := r.db.DB.PrepareContext(ctx, InsertContact)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, params.PhoneNumber, params.Email, params.LinkedId, params.LinkPrecedence)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func NewContactRepo(db *models.Database) *ContactRepo {
	return &ContactRepo{
		db,
	}
}
