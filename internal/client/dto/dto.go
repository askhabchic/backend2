package dto

import "github.com/gofrs/uuid"

type ClientDTO struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"client_name"`
	Surname          string    `json:"client_surname"`
	Birthday         string    `json:"birthday"`
	Gender           bool      `json:"gender"`
	RegistrationDate string    `json:"registration_date"`
	AddressId        uuid.UUID `json:"address_id"`
}

var ClientTableQuery = `CREATE TABLE IF NOT EXISTS client
(
id uuid NOT NULL,
client_name character varying(20) COLLATE pg_catalog."default",
client_surname character varying(20) COLLATE pg_catalog."default",
birthday character varying(20),
gender boolean,
registration_date character varying(20),
address_id uuid,
CONSTRAINT client_pkey PRIMARY KEY (id),
CONSTRAINT address_id FOREIGN KEY (address_id)
REFERENCES address (id) MATCH SIMPLE
ON UPDATE NO ACTION
ON DELETE NO ACTION
);`
