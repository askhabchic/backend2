package model

import "github.com/gofrs/uuid"

type Address struct {
	ID      uuid.UUID `json:"id"`
	Country string    `json:"country"`
	City    string    `json:"city"`
	Street  string    `json:"street"`
}

var AddressTableQuery = `CREATE TABLE IF NOT EXISTS address
(
    id uuid NOT NULL,
    country character varying(100) COLLATE pg_catalog."default",
    city character varying(100) COLLATE pg_catalog."default",
    street character varying(200) COLLATE pg_catalog."default",
    CONSTRAINT address_pkey PRIMARY KEY (id)
);`

var AddressInsertionQuery = `INSERT INTO address (id, country, city, street) VALUES ($1, $2, $3, $4)`
