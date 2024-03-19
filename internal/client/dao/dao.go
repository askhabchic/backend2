package dao

import (
	model2 "backend2/internal/address/model"
	"backend2/internal/client/db"
	"backend2/internal/client/model"
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

func (c *DAO) Create(ctx context.Context, cl *model.Client) (*model.Client, error) {
	client, err := c.repo.Create(ctx, cl)
	if err != nil {
		return &model.Client{}, err
	}
	return client, nil
}

func (c *DAO) FindOne(ctx context.Context, name, surname string) (*model.Client, error) {
	one, err := c.repo.FindOne(ctx, name, surname)
	if err != nil {
		return &model.Client{}, err
	}
	return one, nil
}

func (c *DAO) FindAll(ctx context.Context, limit, offset int) ([]model.Client, error) {
	all, err := c.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return []model.Client{}, err
	}
	return all, nil
}

func (c *DAO) Update(ctx context.Context, id string, addr model2.Address) error {
	err := c.repo.Update(ctx, id, addr)
	if err != nil {
		return err
	}
	return nil
}

func (c *DAO) Delete(ctx context.Context, id string) error {
	err := c.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
