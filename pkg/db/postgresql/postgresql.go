package postgresql

import (
	"backend2/internal/config"
	"backend2/pkg/logging"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return nil
	}
	return
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username,
		sc.Password, sc.Host, sc.Port, sc.Database)
	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			fmt.Println("failed to connect to postgresql")
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}
	return pool, nil
}

func CreateTable(query string, ctx context.Context, pool *pgxpool.Pool, logger *logging.Logger) error {
	_, err := pool.Exec(ctx, query)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return err
	}
	return nil
}
