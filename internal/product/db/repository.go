package db

import (
	"backend2/internal/product/dao"
	"backend2/internal/product/dto"
	"context"
)

type ProductRepository struct {
	dao *dao.ProductDAO
}

func NewProductRepository(dao *dao.ProductDAO) *ProductRepository {
	return &ProductRepository{dao: dao}
}

func (repo *ProductRepository) Create(ctx context.Context, product *dto.ProductDTO) error {
	return repo.dao.Create(ctx, product)
}

func (repo *ProductRepository) FindOne(ctx context.Context, id string) (*dto.ProductDTO, error) {
	return repo.dao.FindOne(ctx, id)
}

func (repo *ProductRepository) FindAll(ctx context.Context) ([]dto.ProductDTO, error) {
	return repo.dao.FindAll(ctx)
}

func (repo *ProductRepository) Update(ctx context.Context, id string, count int) error {
	return repo.dao.Update(ctx, id, count)
}

func (repo *ProductRepository) Delete(ctx context.Context, id string) error {
	return repo.dao.Delete(ctx, id)
}
