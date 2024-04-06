package models

import "time"

type Contact struct {
	Id             int        `db:"id"`
	PhoneNumber    *string    `db:"phonenumber"`
	Email          *string    `db:"email"`
	LinkedId       *int       `db:"linkedid"`
	LinkPrecedence string     `db:"linkprecedence"`
	CreatedAt      *time.Time `db:"createdat"`
	UpdatedAt      *time.Time `db:"updatedat"`
	DeletedAt      *time.Time `db:"deletedat"`
}

type IdentifyContactReq struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

type GetContactParams struct {
	Email       *string
	PhoneNumber *string
}

type InsertContactParams struct {
	Email          *string
	PhoneNumber    *string
	LinkedId       *int
	LinkPrecedence string
}

type IdentifyContactResp struct {
	Contact IdentifyContact `json:"contact"`
}

type IdentifyContact struct {
	PrimaryContactId    int      `json:"primaryContactId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIds []int    `json:"secondaryContactIds"`
}
