package main

import (
	"backend2/cmd/server"
	addrDao "backend2/internal/address/dao"
	addrDb "backend2/internal/address/db"
	modelAddress "backend2/internal/address/dto"
	addrhandler "backend2/internal/address/handler"
	clientDao "backend2/internal/client/dao"
	cliDb "backend2/internal/client/db"
	modelClient "backend2/internal/client/dto"
	"backend2/internal/client/handler"
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

	logger.Infof("create AddressDTO table")
	err = postgresql.CreateTable(modelAddress.AddressTableQuery, ctx, cli, logger)
	if err != nil {
		return
	}

	logger.Infof("create AddressDTO Service")
	addrDAO := addrDao.NewAddressDAO(cli, logger)
	addRepo := addrDb.NewAddressRepository(addrDAO)

	logger.Infof("create AddressDTO Handler")
	addrHandler := addrhandler.NewAddressHandler(logger, addRepo, ctx)
	addrHandler.Register(router)

	logger.Infof("create ClientDTO table")
	err = postgresql.CreateTable(modelClient.ClientTableQuery, ctx, cli, logger)
	if err != nil {
		return
	}

	logger.Infof("create ClientDTO Service")
	clientDAO := clientDao.NewClientDAO(cli, logger)
	logger.Infof("create ClientDTO Repository")
	clientRepo := cliDb.NewClientRepository(clientDAO)

	logger.Infof("create ClientDTO Handler")
	clientHandler := handler.NewClientHandler(logger, clientRepo, ctx)
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
