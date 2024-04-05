package models

import "time"

type Contact struct {
	Id             int        `db:"id"`
	PhoneNumber    *string    `db:"phoneNumber"`
	Email          *string    `db:"email"`
	LinkedId       *int       `db:"linkedId"`
	LinkPrecedence string     `db:"linkPrecedence"`
	CreatedAt      *time.Time `db:"createdAt"`
	UpdatedAt      *time.Time `db:"updatedAt"`
	DeletedAt      *time.Time `db:"deletedAt"`
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
