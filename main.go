package main

import (
	"os"

	"github.com/huynhtruongson/simple-todo/server"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load configuration")
	}
	if config.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	server := server.NewServer(
		server.WithConfig(config),
		server.WithDB(),
		server.WithTokenMaker(),
		server.WithAuthService(),
		server.WithUserService(),
		server.WithTaskService(),
		server.WithSwaggerDoc(),
	)
	server.RunMigration("file://migration")
	go server.RunTaskWorker()
	// go server.RunGatewayServer()
	go server.RunGRPCServer()
	server.RunGinServer()
}
