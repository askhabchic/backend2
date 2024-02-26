package main

import (
	"backend2/cmd/server"
	model2 "backend2/internal/address/model"
	"backend2/internal/client/dao"
	"backend2/internal/client/db"
	"backend2/internal/client/handler"
	"backend2/internal/client/model"
	"backend2/internal/config"
	"backend2/pkg/db/postgresql"
	"backend2/pkg/logging"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"time"
)

func main() {
	logger := logging.GetLogger()

	logger.Infof("create router")
	router := httprouter.New()

	logger.Infof("get config")
	cfg := config.GetConfig()

	logger.Infof("create server")
	newServer := server.NewServer()

	//handlers
	logger.Infof("create requests handler")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Infof("connect ot database")
	cli, repo, err := connection(ctx, cfg, logger)
	if err != nil {
		logger.Errorf("Error: %v", err)
	}
	defer cli.Close()

	logger.Infof("create Address table")
	err = postgresql.CreateTable(model2.AddressTableQuery, ctx, cli, logger)
	if err != nil {
		return
	}
	logger.Infof("create Client table")
	err = postgresql.CreateTable(model.ClientTableQuery, ctx, cli, logger)
	if err != nil {
		return
	}

	logger.Infof("create Client Service")
	clientDAO := dao.NewClientDAO(repo)
	//service := client.NewClientService(repo)

	logger.Infof("create Client Handler")
	handler := handler.NewClientHandler(logger, clientDAO)
	handler.Register(router)

	logger.Infof("start server")
	newServer.Start(cfg, handler, ctx)
}

func connection(ctx context.Context, cfg *config.Config, logger *logging.Logger) (*pgxpool.Pool, *db.Repository, error) {
	cli, err := postgresql.NewClient(ctx, 5, cfg.Storage)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return nil, nil, err
	}
	repository := db.NewRepository(cli, logger)
	logger.Infof("connection to %s is established", cfg.Storage.Database)
	return cli, repository, nil
}
