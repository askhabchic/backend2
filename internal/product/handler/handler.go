package handler

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

var ProductTableQuery = ``
