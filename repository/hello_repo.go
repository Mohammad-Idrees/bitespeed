package repository

import (
	"bitespeed/models"
	"context"
)

type HelloRepo struct {
	db *models.Database
}

func NewHelloRepo(db *models.Database) *HelloRepo {
	return &HelloRepo{
		db,
	}
}

func (r *HelloRepo) Ping(ctx context.Context) error {
	var result int
	err := r.db.DB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return err
	}
	return nil
}
