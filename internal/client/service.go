package client

import (
	"context"
)

type Service struct {
	r repository
}

func NewClientService(r *Repository) *Service {
	return &Service{
		r,
	}
}

func (c *Service) Create(ctx context.Context, cl *Client) (*Client, error) {
	cli, err := c.r.Create(ctx, cl)
	if err != nil {
		return &Client{}, err
	}
	return cli, nil
}

func (c *Service) FindOne(ctx context.Context, name, surname string) (*Client, error) {
	one, err := c.r.FindOne(ctx, name, surname)
	if err != nil {
		return &Client{}, err
	}
	return one, nil
}

func (c *Service) FindAll(ctx context.Context, limit, offset int) ([]Client, error) {
	all, err := c.r.FindAll(ctx, limit, offset)
	if err != nil {
		return []Client{}, err
	}
	return all, nil
}

func (c *Service) Update(ctx context.Context, id, addr string) error {
	err := c.r.Update(ctx, id, addr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Service) Delete(ctx context.Context, id string) error {
	err := c.r.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
