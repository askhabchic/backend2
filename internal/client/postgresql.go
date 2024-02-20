package client

import (
	"backend2/pkg/logging"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository interface {
	Create(ctx context.Context, cl *Client) (*Client, error)
	FindOne(ctx context.Context, name, surname string) (*Client, error)
	FindAll(ctx context.Context, limit, offset int) ([]Client, error)
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

func (r *Repository) Create(ctx context.Context, cl *Client) (*Client, error) {
	q := `INSERT INTO client (client_name, client_surname, birthday, gender, registration_date, address_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	err := r.psgr.QueryRow(ctx, q, cl.Name, cl.Surname, cl.Birthday, cl.Gender, cl.RegistrationDate, cl.AddressId).Scan(&cl.ID)
	if err != nil {
		return &Client{}, err
	}
	return cl, nil
}

func (r *Repository) FindAll(ctx context.Context, limit, offset int) (cls []Client, err error) {
	q := `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM public.client`
	if limit != 0 {
		q = fmt.Sprintf(q + ` LIMIT $1`)
	}
	if offset != 0 {
		q = fmt.Sprintf(q + ` OFFSET $2`)
	}
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	rows, err := r.psgr.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	cls = make([]Client, 0)
	for rows.Next() {
		var cl Client
		err := rows.Scan(&cl.ID, &cl.Name, &cl.Surname, &cl.Birthday, &cl.Gender, &cl.RegistrationDate, &cl.AddressId)
		if err != nil {
			return nil, err
		}
		cls = append(cls, cl)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cls, nil
}

func (r *Repository) FindOne(ctx context.Context, name, surname string) (*Client, error) {
	q := `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM public.client WHERE client_name = $1 client_surname = $2`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	var cl Client
	if err := r.psgr.QueryRow(ctx, q, name, surname).Scan(&cl.ID, &cl.Name, &cl.Surname, &cl.Birthday, &cl.Gender, &cl.RegistrationDate, &cl.AddressId); err != nil {
		return &Client{}, err
	}
	return &cl, nil
}

func (r *Repository) Update(ctx context.Context, id, addr string) error {
	q := `UPDATE client SET address_id = $1 WHERE id = $2`

	rows, err := r.psgr.Query(ctx, q, addr, id)
	if err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM client WHERE id = $1`

	rows, err := r.psgr.Query(ctx, q, id)
	if err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
