package dto

import (
	"github.com/gofrs/uuid"
	"time"
)

type ProductDTO struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Price          float64   `json:"price"`
	AvailableStock int       `json:"available_stock"`
	LastUpdateDate time.Time `json:"last_update_date"`
	SupplierId     uuid.UUID `json:"supplier_id"`
	ImageId        uuid.UUID `json:"image_id"`
}

var ProductTableQuery = `CREATE TABLE IF NOT EXISTS product
(
    id uuid NOT NULL,
    name character varying(200) COLLATE pg_catalog."default",
    category character varying(100) COLLATE pg_catalog."default",
    price money,
    available_stock integer,
    last_update_date date,
    supplier_id uuid,
    image_id uuid,
    CONSTRAINT product_pkey PRIMARY KEY (id),
    CONSTRAINT image_id FOREIGN KEY (image_id)
        REFERENCES images (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT supplier_id FOREIGN KEY (supplier_id)
        REFERENCES supplier (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);`
