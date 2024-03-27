package db

import (
	"backend2/internal/image/dao"
	"backend2/internal/image/dto"
	"context"
)

type ImageRepository struct {
	dao *dao.ImageDAO
}

func NewImageRepository(dao *dao.ImageDAO) *ImageRepository {
	return &ImageRepository{
		dao: dao,
	}
}

func (repo *ImageRepository) Create(ctx context.Context, imageBody []byte, id string) error {
	return repo.dao.Create(ctx, imageBody, id)
}

func (repo *ImageRepository) FindOne(ctx context.Context, id string) (*dto.ImageDTO, error) {
	return repo.dao.FindOne(ctx, id)
}

func (repo *ImageRepository) FindProductImage(ctx context.Context, productID string) (*dto.ImageDTO, error) {
	return repo.dao.FindOne(ctx, productID)
}

func (repo *ImageRepository) Update(ctx context.Context, id, newImage string) error {
	return repo.dao.Update(ctx, id, newImage)
}

func (repo *ImageRepository) Delete(ctx context.Context, id string) error {
	return repo.dao.Delete(ctx, id)
}
