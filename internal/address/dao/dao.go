package dao

import (
	"backend2/internal/address/db"
	"backend2/internal/address/model"
	"context"
)

type DAO struct {
	repo *db.Repository
}

func NewClientDAO(r *db.Repository) *DAO {
	return &DAO{
		r,
	}
}

func (c *DAO) Create(ctx context.Context, cl *model.Address) (*model.Address, error) {
	cli, err := c.repo.Create(ctx, cl)
	if err != nil {
		return &model.Address{}, err
	}
	return cli, nil
}

func (c *DAO) FindOne(ctx context.Context, name, surname string) (*model.Address, error) {
	one, err := c.repo.FindOne(ctx, name, surname)
	if err != nil {
		return &model.Address{}, err
	}
	return one, nil
}

func (c *DAO) FindAll(ctx context.Context, limit, offset int) ([]model.Address, error) {
	all, err := c.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return []model.Address{}, err
	}
	return all, nil
}

func (c *DAO) Update(ctx context.Context, id string, addr *model.Address) (*model.Address, error) {
	updatedAddress, err := c.repo.Update(ctx, id, addr)
	if err != nil {
		return nil, err
	}
	return updatedAddress, nil
}

func (c *DAO) Delete(ctx context.Context, id string) error {
	err := c.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
