package dao

import (
	addressDAO "backend2/internal/address/dao"
	model2 "backend2/internal/address/dto"
	"backend2/internal/client/dto"
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

func (dao *ClientDAO) Create(ctx context.Context, cl *dto.ClientDTO) error {
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

func (dao *ClientDAO) FindAll(ctx context.Context, limit, offset string) (cls []dto.ClientDTO, err error) {
	q := `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM public.client`
	if limit != "" {
		q = fmt.Sprintf(q + fmt.Sprintf(` LIMIT %s`, limit))
	}
	if offset != "" {
		q = fmt.Sprintf(q + fmt.Sprintf(` OFFSET %s`, offset))
	}
	dao.logger.Tracef("SQL Query: %s", q)
	rows, err := dao.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}
	cls = make([]dto.ClientDTO, 0)
	for rows.Next() {
		var cl dto.ClientDTO
		err = rows.Scan(&cl.ID, &cl.Name, &cl.Surname, &cl.Birthday,
			&cl.Gender, &cl.RegistrationDate, &cl.AddressId)
		if err != nil && errors.Is(err, rows.Err()) {
			return nil, err
		}
		cls = append(cls, cl)
	}
	return cls, nil
}

func (dao *ClientDAO) FindOne(ctx context.Context, name, surname string) (*dto.ClientDTO, error) {
	q := `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM public.client WHERE client_name=$1 AND client_surname = $2`

	dao.logger.Tracef("SQL Query: %s", q)
	var cl dto.ClientDTO
	if err := dao.db.QueryRow(ctx, q, name, surname).Scan(&cl.ID, &cl.Name, &cl.Surname,
		&cl.Birthday, &cl.Gender, &cl.RegistrationDate, &cl.AddressId); err != nil {
		return &dto.ClientDTO{}, err
	}
	return &cl, nil
}

func (dao *ClientDAO) Update(ctx context.Context, id string, addr *model2.AddressDTO) error {
	var existedID string
	querySelect := `SELECT id FROM address WHERE city = $1 AND street = $2`
	dao.logger.Tracef("SQL Query: %s", model2.AddressInsertionQuery)
	err := dao.db.QueryRow(ctx, querySelect, addr.City, addr.Street).Scan(&existedID)

	queryUpdate := `UPDATE client SET address_id = $1 WHERE id = $2`
	dao.logger.Tracef("SQL Query: %s", queryUpdate)

	if err != nil && err.Error() == pgx.ErrNoRows.Error() {
		newAddress, err := addressDAO.NewAddressDAO(dao.db, dao.logger).Create(ctx, addr)
		if err != nil {
			return err
		}

		_, queryError := dao.db.Query(ctx, queryUpdate, newAddress.ID, id)
		if queryError != nil {
			return queryError
		}
		return nil
	}

	_, queryError := dao.db.Query(ctx, queryUpdate, existedID, id)
	if queryError != nil {
		return queryError
	}

	return nil
}

func (dao *ClientDAO) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM client WHERE id = $1`
	dao.logger.Tracef("SQL Query: %s", q)

	_, err := dao.db.Query(ctx, q, id)
	if err != nil {
		return err
	}
	return nil
}
