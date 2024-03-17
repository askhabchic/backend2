package db

import (
	"backend2/internal/address/model"
	"backend2/pkg/logging"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

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

func (r Repository) Create(ctx context.Context, addr *model.Address) (*model.Address, error) {
	q := `INSERT INTO address (id, country, city, street) VALUES ($1, $2, $3, $4) RETURNING id`
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	r.logger.Trace(fmt.Sprint("SQL Query: %s", q))
	err = r.psgr.QueryRow(ctx, q, id, addr.Country, addr.City, addr.Street).Scan(&addr.ID)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (r Repository) FindOne(ctx context.Context, id string) (*model.Address, error) {
	q := `SELECT id, country, city, streer FROM public.address WHERE id = $1`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	var addr model.Address
	if err := r.psgr.QueryRow(ctx, q, id).Scan(&addr.ID, &addr.Country, &addr.City, &addr.Street); err != nil {
		return &model.Address{}, err
	}
	return &addr, nil
}

func (r Repository) FindAll(ctx context.Context) ([]model.Address, error) {
	q := `SELECT id, country, city, street FROM public.address`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	rows, err := r.psgr.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	addrs := make([]model.Address, 0)
	for rows.Next() {
		var addr model.Address
		err := rows.Scan(&addr.ID, &addr.Country, &addr.City, &addr.Street)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return addrs, nil
}

func (r Repository) Update(ctx context.Context, id string, addr *model.Address) (*model.Address, error) {
	q := `UPDATE address SET country = $1, city = $2, street = $3 WHERE id = $4`

	rows, _ := r.psgr.Query(ctx, q, addr.Country, addr.City, addr.Street, id)
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return addr, nil
}

func (r Repository) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM address WHERE id = $1`

	rows, err := r.psgr.Query(ctx, q, id)
	if err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
