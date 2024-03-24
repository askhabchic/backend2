package db

import (
	"backend2/internal/address/dao"
	"backend2/internal/address/dto"
	"context"
)

type AddressRepository struct {
	dao *dao.AddressDAO
}

func NewAddressRepository(dao *dao.AddressDAO) *AddressRepository {
	return &AddressRepository{
		dao,
	}
}

func (repo *AddressRepository) Create(ctx context.Context, cl *dto.AddressDTO) (*dto.AddressDTO, error) {
	return repo.dao.Create(ctx, cl)
}

func (repo *AddressRepository) FindOne(ctx context.Context, id string) (*dto.AddressDTO, error) {
	return repo.dao.FindOne(ctx, id)
}

func (repo *AddressRepository) FindAll(ctx context.Context) ([]dto.AddressDTO, error) {
	return repo.dao.FindAll(ctx)
}

func (repo *AddressRepository) Update(ctx context.Context, id string, addr *dto.AddressDTO) (*dto.AddressDTO, error) {
	return repo.dao.Update(ctx, id, addr)
}

func (repo *AddressRepository) Delete(ctx context.Context, id string) error {
	return repo.dao.Delete(ctx, id)
}
