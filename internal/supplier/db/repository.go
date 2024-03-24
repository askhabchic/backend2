package db

import (
	model2 "backend2/internal/address/dto"
	"backend2/internal/supplier/dao"
	"backend2/internal/supplier/dto"
	"context"
)

type SupplierRepository struct {
	dao *dao.SupplierDAO
}

func NewSupplierRepository(dao *dao.SupplierDAO) *SupplierRepository {
	return &SupplierRepository{
		dao: dao,
	}
}

func (repo *SupplierRepository) Create(ctx context.Context, supplier *dto.SupplierDTO) error {
	return repo.dao.Create(ctx, supplier)
}

func (repo *SupplierRepository) FindOne(ctx context.Context, id string) (*dto.SupplierDTO, error) {
	return repo.dao.FindOne(ctx, id)
}

func (repo *SupplierRepository) FindAll(ctx context.Context) ([]dto.SupplierDTO, error) {
	return repo.dao.FindAll(ctx)
}

func (repo *SupplierRepository) Update(ctx context.Context, id string, addr model2.AddressDTO) error {
	return repo.dao.Update(ctx, id, &addr)
}

func (repo *SupplierRepository) Delete(ctx context.Context, id string) error {
	return repo.dao.Delete(ctx, id)
}
