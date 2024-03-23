package dao

import (
	model2 "backend2/internal/address/dto"
	"backend2/internal/supplier/dto"
	"backend2/pkg/logging"
	"context"
	"github.com/gofrs/uuid"
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
	return nil, nil
}

func (dao *SupplierDAO) FindAll(ctx context.Context) ([]dto.SupplierDTO, error) {
	return nil, nil
}

func (dao *SupplierDAO) Update(ctx context.Context, id string, address model2.AddressDTO) error {
	return nil
}

func (dao *SupplierDAO) Delete(ctx context.Context, id string) error {
	return nil
}
