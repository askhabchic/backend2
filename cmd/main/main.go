package main

import (
	"backend2/cmd/server"
	addrDao "backend2/internal/address/dao"
	addrDb "backend2/internal/address/db"
	addrModel "backend2/internal/address/dto"
	addrhandler "backend2/internal/address/handler"
	clientDao "backend2/internal/client/dao"
	cliDb "backend2/internal/client/db"
	clientModel "backend2/internal/client/dto"
	clientHandler "backend2/internal/client/handler"
	"backend2/internal/config"
	supplierDAO "backend2/internal/supplier/dao"
	supplierDb "backend2/internal/supplier/db"
	supplierModel "backend2/internal/supplier/dto"
	supplierHandler "backend2/internal/supplier/handler"
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

	//context creation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Infof("connect ot database")
	cli, err := connection(ctx, cfg, logger)
	if err != nil {
		logger.Errorf("Error: %v", err)
	}
	defer cli.Close()

	err = createTables(ctx, cli, logger)
	if err != nil {
		return
	}

	initializeRouter(logger, cli, ctx, router)

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

func createTables(ctx context.Context, cli *pgxpool.Pool, logger *logging.Logger) error {
	logger.Infof("create AddressDTO table")
	err := postgresql.CreateTable(addrModel.AddressTableQuery, ctx, cli, logger)
	if err != nil {
		return err
	}

	logger.Infof("create ClientDTO table")
	err = postgresql.CreateTable(clientModel.ClientTableQuery, ctx, cli, logger)
	if err != nil {
		return err
	}

	logger.Infof("create Supplier table")
	err = postgresql.CreateTable(supplierModel.SupplierTableQuery, ctx, cli, logger)
	if err != nil {
		return err
	}
	return nil
}

func initializeRouter(logger *logging.Logger, cli *pgxpool.Pool, ctx context.Context, router *httprouter.Router) {
	logger.Infof("create Address Service")
	addrDAO := addrDao.NewAddressDAO(cli, logger)
	logger.Infof("create Address Repository")
	addRepo := addrDb.NewAddressRepository(addrDAO)
	logger.Infof("create Address Handler")
	addrHandler := addrhandler.NewAddressHandler(logger, addRepo, ctx)
	addrHandler.Register(router)

	logger.Infof("create Client Service")
	clientDAO := clientDao.NewClientDAO(cli, logger)
	logger.Infof("create Client Repository")
	clientRepo := cliDb.NewClientRepository(clientDAO)
	logger.Infof("create Client Handler")
	clHandler := clientHandler.NewClientHandler(logger, clientRepo, ctx)
	clHandler.Register(router)

	logger.Infof("create Supplier Service")
	supplierDAO := supplierDAO.NewSupplierDAO(cli, logger)
	logger.Infof("create Supplier Repository")
	supplierRepo := supplierDb.NewSupplierRepository(supplierDAO)
	logger.Infof("create Supplier Handler")
	supplHandler := supplierHandler.NewSupplierHandler(logger, supplierRepo, ctx)
	supplHandler.Register(router)
}
