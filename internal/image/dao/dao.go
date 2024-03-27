package dao

import (
	"backend2/internal/image/dto"
	"backend2/pkg/logging"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ImageDAO struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewImageDAO(client *pgxpool.Pool, logger *logging.Logger) *ImageDAO {
	return &ImageDAO{
		db:     client,
		logger: logger,
	}
}

func (dao *ImageDAO) Create(ctx context.Context, imageBody []byte, id string) error {
	return nil
}

func (dao *ImageDAO) FindOne(ctx context.Context, id string) (*dto.ImageDTO, error) {
	return nil, nil
}

func (dao *ImageDAO) FindProductImage(ctx context.Context, productID string) (*dto.ImageDTO, error) {
	return nil, nil
}

func (dao *ImageDAO) Update(ctx context.Context, id, newImage string) error {
	return nil
}

func (dao *ImageDAO) Delete(ctx context.Context, id string) error {
	return nil
}
