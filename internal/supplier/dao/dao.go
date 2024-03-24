package dao

import (
	addressDAO "backend2/internal/address/dao"
	model2 "backend2/internal/address/dto"
	"backend2/internal/supplier/dto"
	"backend2/pkg/logging"
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SupplierDAO struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewSupplierDAO(client *pgxpool.Pool, logger *logging.Logger) *SupplierDAO {
	return &SupplierDAO{
		db:     client,
		logger: logger,
	}
}

func (dao *SupplierDAO) Create(ctx context.Context, sup *dto.SupplierDTO) error {
	q := `INSERT INTO supplier (id, name, address_id, phone_number) VALUES ($1, $2, $3, $4) RETURNING id`

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	dao.logger.Tracef("SQL Query: %s", q)
	err = dao.db.QueryRow(ctx, q, id, sup.Name, sup.AddressId, sup.PhoneNumber).Scan(&sup.ID)
	if err != nil {
		return err
	}
	return nil
}

func (dao *SupplierDAO) FindOne(ctx context.Context, id string) (*dto.SupplierDTO, error) {
	q := `SELECT id, name, address_id, phone_number FROM supplier WHERE id = $1`
	dao.logger.Tracef("SQL Query: %s", q)

	var sup dto.SupplierDTO
	err := dao.db.QueryRow(ctx, q, id).Scan(&sup.ID, &sup.Name, &sup.AddressId, &sup.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return &sup, nil
}

func (dao *SupplierDAO) FindAll(ctx context.Context) ([]dto.SupplierDTO, error) {
	q := `SELECT id, name, address_id, phone_number FROM supplier`
	dao.logger.Tracef("SQL Query: %s", q)

	rows, err := dao.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var supls []dto.SupplierDTO
	for rows.Next() {
		var sup dto.SupplierDTO
		err = rows.Scan(&sup.ID, &sup.Name, &sup.AddressId, &sup.PhoneNumber)
		if err != nil {
			return nil, err
		}
		supls = append(supls, sup)
	}

	return supls, nil
}

func (dao *SupplierDAO) Update(ctx context.Context, id string, address *model2.AddressDTO) error {
	var existedID string
	querySelect := `SELECT id FROM address WHERE city = $1 AND street = $2`
	dao.logger.Tracef("SQL Query: %s", querySelect)
	err := dao.db.QueryRow(ctx, querySelect, address.City, address.Street).Scan(&existedID)

	queryUpdate := `UPDATE supplier SET address_id = $1 WHERE id = $2`
	dao.logger.Tracef("SQL Query: %s", queryUpdate)

	if err != nil && err.Error() == pgx.ErrNoRows.Error() {
		newAddress, err := addressDAO.NewAddressDAO(dao.db, dao.logger).Create(ctx, address)
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

func (dao *SupplierDAO) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM supplier WHERE id = $1`
	dao.logger.Tracef("SQL query: %s", q)

	rows, err := dao.db.Query(ctx, q, id)
	if err != nil && errors.Is(err, rows.Err()) {
		return err
	}
	return nil
}
