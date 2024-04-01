package models

import (
	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB *sqlx.DB
}
