package dao

import (
	"backend2/internal/product/dto"
	"backend2/pkg/logging"
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductDAO struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

func NewProductDAO(client *pgxpool.Pool, logger logging.Logger) *ProductDAO {
	return &ProductDAO{
		db:     client,
		logger: logger,
	}
}

func (dao *ProductDAO) Create(ctx context.Context, prod *dto.ProductDTO) error {
	q := `INSERT INTO product (id, name, category, price, available_stock, last_update_date, supplier_id, image_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	dao.logger.Tracef("SQL query: %s", q)
	err = dao.db.QueryRow(ctx, q, id, prod.Name, prod.Category, prod.Price,
		prod.AvailableStock, prod.LastUpdateDate, prod.SupplierId, prod.ImageId).Scan(&prod.ID)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ProductDAO) FindOne(ctx context.Context, id string) (*dto.ProductDTO, error) {
	q := `SELECT id, name, category, price, available_stock, last_update_date, supplier_id, image_id FROM product WHERE id = $1`

	dao.logger.Tracef("SQL query: %s", q)
	var product dto.ProductDTO
	err := dao.db.QueryRow(ctx, q, id).Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.AvailableStock,
		&product.LastUpdateDate, &product.SupplierId, &product.ImageId)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (dao *ProductDAO) FindAll(ctx context.Context) ([]dto.ProductDTO, error) {
	q := `SELECT id, name, category, price, available_stock, last_update_date, supplier_id, image_id FROM product`

	dao.logger.Tracef("SQL query: %s", q)
	rows, err := dao.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var products []dto.ProductDTO
	for rows.Next() {
		var product dto.ProductDTO
		err = rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.AvailableStock,
			&product.LastUpdateDate, &product.SupplierId, &product.ImageId)
		if err != nil && errors.Is(err, rows.Err()) {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (dao *ProductDAO) Update(ctx context.Context, id string, count int) error {
	querySelect := `SELECT available_stock FROM product WHERE id = $1`
	queryUpdate := `UPDATE product SET available_stock = $1 WHERE id = $2`

	var amount int
	err := dao.db.QueryRow(ctx, querySelect, id).Scan(&amount)
	if err != nil {
		return err
	}
	amount -= count
	if amount < 0 {
		amount = 0
	}
	_, queryError := dao.db.Query(ctx, queryUpdate, amount, id)
	if queryError != nil {
		return queryError
	}
	return nil
}

func (dao *ProductDAO) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM product WHERE id = $1`
	dao.logger.Tracef("SQL query: %s", q)

	_, err := dao.db.Query(ctx, q, id)
	if err != nil {
		return err
	}
	return nil
}
