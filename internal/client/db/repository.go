package db

import (
	model2 "backend2/internal/address/dto"
	"backend2/internal/client/dao"
	"backend2/internal/client/dto"
	"context"
)

type ClientRepository struct {
	dao *dao.ClientDAO
}

func NewClientRepository(dao *dao.ClientDAO) *ClientRepository {
	return &ClientRepository{
		dao,
	}
}

func (repo *ClientRepository) Create(ctx context.Context, cl *dto.ClientDTO) error {
	return repo.dao.Create(ctx, cl)
}

func (repo *ClientRepository) FindOne(ctx context.Context, name, surname string) (*dto.ClientDTO, error) {
	one, err := repo.dao.FindOne(ctx, name, surname)
	if err != nil {
		return &dto.ClientDTO{}, err
	}
	return one, nil
}

func (repo *ClientRepository) FindAll(ctx context.Context, limit, offset int) ([]dto.ClientDTO, error) {
	all, err := repo.dao.FindAll(ctx, limit, offset)
	if err != nil {
		return []dto.ClientDTO{}, err
	}
	return all, nil
}

func (repo *ClientRepository) Update(ctx context.Context, id string, addr model2.AddressDTO) error {
	return repo.dao.Update(ctx, id, addr)
}

func (repo *ClientRepository) Delete(ctx context.Context, id string) error {
	return repo.dao.Delete(ctx, id)
}
