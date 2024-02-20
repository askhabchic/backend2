package main

import (
	"backend2/cmd/server"
	"backend2/internal/client"
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
	//TODO create context for handlers
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//handl :=

	logger.Infof("connect ot database")
	cli, repo, err := connection(ctx, cfg, logger)
	if err != nil {
		logger.Errorf("Error: %v", err)
	}
	defer cli.Close()

	logger.Infof("create Client Service")
	service := client.NewClientService(repo)

	logger.Infof("create Client Handler")
	handler := client.NewClientHandler(logger, service)
	handler.Register(router)

	logger.Infof("start server")
	//router - handler.Register
	newServer.Start(cfg, handler, ctx)
}

func connection(ctx context.Context, cfg *config.Config, logger *logging.Logger) (*pgxpool.Pool, *client.Repository, error) {
	cli, err := postgresql.NewClient(ctx, 5, cfg.Storage)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return nil, nil, err
	}
	repository := client.NewRepository(cli, logger)
	logger.Infof("connection to %s is established", cfg.Storage.Database)
	return cli, repository, nil
}
