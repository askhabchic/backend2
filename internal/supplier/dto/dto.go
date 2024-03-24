package dto

import "github.com/gofrs/uuid"

type SupplierDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	AddressId   uuid.UUID `json:"address_id"`
	PhoneNumber string    `json:"phone_number"`
}

var SupplierTableQuery = `CREATE TABLE IF NOT EXISTS supplier
(
    id uuid NOT NULL,
    name character varying(100) COLLATE pg_catalog."default",
    address_id uuid,
    phone_number character varying(20) COLLATE pg_catalog."default",
    CONSTRAINT supplier_pkey PRIMARY KEY (id),
    CONSTRAINT address_id FOREIGN KEY (address_id)
    REFERENCES address (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    );`
