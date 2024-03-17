package main

import (
	"backend2/cmd/server"
	addrDao "backend2/internal/address/dao"
	addrDb "backend2/internal/address/db"
	addrhandler "backend2/internal/address/handler"
	modelAddress "backend2/internal/address/model"
	clientDao "backend2/internal/client/dao"
	cliDb "backend2/internal/client/db"
	"backend2/internal/client/handler"
	modelClient "backend2/internal/client/model"
	"backend2/internal/config"
	"backend2/pkg/db/postgresql"
	"backend2/pkg/logging"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logging.GetLogger()

	logger.Infof("create router")
	router := httprouter.New()

	logger.Infof("get config")
	cfg := config.GetConfig()

	logger.Infof("create server")
	newServer := server.NewServer()

	//channel for signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//handlers
	logger.Infof("create requests handler")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Infof("connect ot database")
	cli, err := connection(ctx, cfg, logger)
	if err != nil {
		logger.Errorf("Error: %v", err)
	}
	defer cli.Close()

	logger.Infof("create Address table")
	err = postgresql.CreateTable(modelAddress.AddressTableQuery, ctx, cli, logger)
	if err != nil {
		return
	}

	logger.Infof("create Address Service")
	addRepo := addrDb.NewRepository(cli, logger)
	addrDAO := addrDao.NewAddressDAO(addRepo)

	logger.Infof("create Address Handler")
	addrHandler := addrhandler.NewAddressHandler(logger, addrDAO, ctx)
	addrHandler.Register(router)

	logger.Infof("create Client table")
	err = postgresql.CreateTable(modelClient.ClientTableQuery, ctx, cli, logger)
	if err != nil {
		return
	}

	logger.Infof("create Client Service")
	clientRepo := cliDb.NewRepository(cli, logger)
	clientDAO := clientDao.NewClientDAO(clientRepo)

	logger.Infof("create Client Handler")
	clientHandler := handler.NewClientHandler(logger, clientDAO, ctx)
	clientHandler.Register(router)

	logger.Infof("start server")
	newServer.Start(cfg, router, ctx)

	sig := <-signalCh
	logger.Infof("received signal: %v/n", sig)
}

func connection(ctx context.Context, cfg *config.Config, logger *logging.Logger) (*pgxpool.Pool, error) {
	cli, err := postgresql.NewClient(ctx, 5, cfg.Storage)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return nil, err
	}
	logger.Infof("connection to %s is established", cfg.Storage.Database)
	return cli, nil
}
