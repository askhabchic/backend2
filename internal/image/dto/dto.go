package dto

import "github.com/gofrs/uuid"

type CreateImageDTO struct {
	ID    uuid.UUID `json:"id"`
	Image string    `json:"image"`
}

var ImageTableuery = `
CREATE TABLE IF NOT EXISTS images
(
    id uuid NOT NULL,
    image bytea,
    CONSTRAINT images_pkey PRIMARY KEY (id)
);`
