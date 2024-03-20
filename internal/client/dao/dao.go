package dao

import (
	model2 "backend2/internal/address/model"
	"backend2/internal/client/model"
	"backend2/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ClientDAO struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewClientDAO(client *pgxpool.Pool, logger *logging.Logger) *ClientDAO {
	return &ClientDAO{
		db:     client,
		logger: logger,
	}
}

func (dao *ClientDAO) Create(ctx context.Context, cl *model.Client) error {
	q := `INSERT INTO client (id, client_name, client_surname, birthday, gender, registration_date, address_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id, err = uuid.NewV4()
	if err != nil {
		return err
	}

	dao.logger.Tracef("SQL Query: %s", q)
	err = dao.db.QueryRow(ctx, q, id, cl.Name, cl.Surname, cl.Birthday,
		cl.Gender, cl.RegistrationDate, cl.AddressId).Scan(&cl.ID)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ClientDAO) FindAll(ctx context.Context, limit, offset int) (cls []model.Client, err error) {
	q := `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM public.client`
	if limit != 0 {
		q = fmt.Sprintf(q+` LIMIT %d`, limit)
	}
	if offset != 0 {
		q = fmt.Sprintf(q+` OFFSET %d`, offset)
	}
	dao.logger.Tracef("SQL Query: %s", q)
	rows, err := dao.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}
	cls = make([]model.Client, 0)
	for rows.Next() {
		var cl model.Client
		err = rows.Scan(&cl.ID, &cl.Name, &cl.Surname, &cl.Birthday,
			&cl.Gender, &cl.RegistrationDate, &cl.AddressId)
		if err != nil && errors.Is(err, rows.Err()) {
			return nil, err
		}
		cls = append(cls, cl)
	}
	return cls, nil
}

func (dao *ClientDAO) FindOne(ctx context.Context, name, surname string) (*model.Client, error) {
	q := `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM public.client WHERE client_name=$1 AND client_surname = $2`

	dao.logger.Tracef("SQL Query: %s", q)
	var cl model.Client
	if err := dao.db.QueryRow(ctx, q, name, surname).Scan(&cl.ID, &cl.Name, &cl.Surname,
		&cl.Birthday, &cl.Gender, &cl.RegistrationDate, &cl.AddressId); err != nil {
		return &model.Client{}, err
	}
	return &cl, nil
}

func (dao *ClientDAO) Update(ctx context.Context, id string, addr model2.Address) error {
	idUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	var existedID uuid.UUID //TODO it's maybe problem
	q := `SELECT id FROM address WHERE city = $1 AND street = $2`
	dao.logger.Tracef("SQL Query: %s", model2.AddressInsertionQuery)
	err = dao.db.QueryRow(ctx, q, addr.City, addr.Street).Scan(&existedID)

	q = `UPDATE client SET address_id = $1 WHERE id = $2`
	dao.logger.Tracef("SQL Query: %s", q)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		row := dao.db.QueryRow(ctx, model2.AddressInsertionQuery, idUUID, addr.Country, addr.City, addr.Street)
		if err = row.Scan(&addr.ID); err != nil {
			return err
		}
		//todo implementation of correct return
		_, queryError := dao.db.Query(ctx, q, idUUID, id)
		if queryError != nil {
			return queryError
		}
	}

	rows, _ := dao.db.Query(ctx, q, existedID, id)
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func (dao *ClientDAO) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM client WHERE id = $1`
	dao.logger.Tracef("SQL Query: %s", q)

	rows, err := dao.db.Query(ctx, q, id)
	if err != nil || errors.Is(err, rows.Err()) {
		return err
	}
	return nil
}
