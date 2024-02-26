package db

import (
	"backend2/internal/address/model"
	"backend2/pkg/logging"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository interface {
	Create(ctx context.Context, cl *model.Address) (*model.Address, error)
	FindOne(ctx context.Context, name, surname string) (*model.Address, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.Address, error)
	Update(ctx context.Context, id, addr string) error
	Delete(ctx context.Context, id string) error
}

type Repository struct {
	psgr   *pgxpool.Pool
	logger *logging.Logger
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		psgr:   client,
		logger: logger,
	}
}

func (r Repository) Create(ctx context.Context, aadr *model.Address) (*model.Address, error) {
	return nil, nil
}

func (r Repository) FindOne(ctx context.Context, name string, surname string) (*model.Address, error) {
	return nil, nil
}

func (r Repository) FindAll(ctx context.Context, limit int, offset int) ([]model.Address, error) {
	return nil, nil
}

func (r Repository) Update(ctx context.Context, id string, addr string) error {
	return nil
}

func (r Repository) Delete(ctx context.Context, id string) error {
	return nil
}
