package model

import "github.com/gofrs/uuid"

type Supplier struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	AddressId   uuid.UUID `json:"address_id"`
	PhoneNumber string    `json:"phone_number"`
}
