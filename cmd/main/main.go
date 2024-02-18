package main

import (
	"backend2/cmd/server"
	"backend2/internal/config"
	"backend2/pkg/db/postgresql"
	"backend2/pkg/logging"
	"context"
	"github.com/julienschmidt/httprouter"
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
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//handl :=

	logger.Infof("start server")
	newServer.Start(cfg, router)

	logger.Infof("connect ot database")
	client, err := postgresql.NewClient(context.Background(), 5, cfg.Storage)
	if err != nil {
		logger.Errorf("Error: %v", err)
	}
	logger.Infof("connection to %s is established", cfg.Storage.Database)

	defer client.Close()
}
