package dao

import (
	"backend2/internal/address/dto"
	"backend2/pkg/logging"
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

//type DAO struct {
//	repo *db.Repository
//}
//
//func NewAddressDAO(r *db.Repository) *DAO {
//	return &DAO{
//		r,
//	}
//}
//
//func (c *DAO) Create(ctx context.Context, cl *dto.AddressDTO) (*dto.AddressDTO, error) {
//	cli, err := c.repo.Create(ctx, cl)
//	if err != nil {
//		return &dto.AddressDTO{}, err
//	}
//	return cli, nil
//}
//
//func (c *DAO) FindOne(ctx context.Context, id string) (*dto.AddressDTO, error) {
//	one, err := c.repo.FindOne(ctx, id)
//	if err != nil {
//		return &dto.AddressDTO{}, err
//	}
//	return one, nil
//}
//
//func (c *DAO) FindAll(ctx context.Context) ([]dto.AddressDTO, error) {
//	all, err := c.repo.FindAll(ctx)
//	if err != nil {
//		return []dto.AddressDTO{}, err
//	}
//	return all, nil
//}
//
//func (c *DAO) Update(ctx context.Context, id string, addr *dto.AddressDTO) (*dto.AddressDTO, error) {
//	updatedAddress, err := c.repo.Update(ctx, id, addr)
//	if err != nil {
//		return nil, err
//	}
//	return updatedAddress, nil
//}
//
//func (c *DAO) Delete(ctx context.Context, id string) error {
//	err := c.repo.Delete(ctx, id)
//	if err != nil {
//		return err
//	}
//	return nil
//}

type AddressDAO struct {
	psgr   *pgxpool.Pool
	logger *logging.Logger
}

func NewAddressDAO(client *pgxpool.Pool, logger *logging.Logger) *AddressDAO {
	return &AddressDAO{
		psgr:   client,
		logger: logger,
	}
}

func (dao *AddressDAO) Create(ctx context.Context, addr *dto.AddressDTO) (*dto.AddressDTO, error) {
	q := `INSERT INTO address (id, country, city, street) VALUES ($1, $2, $3, $4) RETURNING id`
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	dao.logger.Tracef("SQL Query: %s", q)
	err = dao.psgr.QueryRow(ctx, q, id, addr.Country, addr.City, addr.Street).Scan(&addr.ID)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (dao *AddressDAO) FindOne(ctx context.Context, id string) (*dto.AddressDTO, error) {
	q := `SELECT id, country, city, streer FROM public.address WHERE id = $1`

	dao.logger.Tracef("SQL Query: %s", q)
	var addr dto.AddressDTO
	if err := dao.psgr.QueryRow(ctx, q, id).Scan(&addr.ID, &addr.Country, &addr.City, &addr.Street); err != nil {
		return &dto.AddressDTO{}, err
	}
	return &addr, nil
}

func (dao *AddressDAO) FindAll(ctx context.Context) ([]dto.AddressDTO, error) {
	q := `SELECT id, country, city, street FROM public.address`

	dao.logger.Tracef("SQL Query: %s", q)
	rows, err := dao.psgr.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	addrs := make([]dto.AddressDTO, 0)
	for rows.Next() {
		var addr dto.AddressDTO
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

func (dao *AddressDAO) Update(ctx context.Context, id string, addr *dto.AddressDTO) (*dto.AddressDTO, error) {
	q := `UPDATE address SET country = $1, city = $2, street = $3 WHERE id = $4`

	rows, _ := dao.psgr.Query(ctx, q, addr.Country, addr.City, addr.Street, id)
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return addr, nil
}

func (dao *AddressDAO) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM address WHERE id = $1`

	rows, err := dao.psgr.Query(ctx, q, id)
	if err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
