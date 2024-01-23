package common

import (
	"context"
	"fmt"

	"github.com/huynhtruongson/simple-todo/lib"
	"github.com/huynhtruongson/simple-todo/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Suite struct {
	DB       lib.DB
	Config   utils.Config
	BaseURL  string
	GrpcConn grpc.ClientConnInterface
}

func NewSuite() *Suite {
	config, err := utils.LoadConfig("../.")
	if err != nil {
		log.Fatal().Err(err).Msg("load config failed")
	}
	db, err := pgxpool.New(context.Background(), config.DBAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
	grpcConn, err := grpc.Dial(config.GRPCServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to grpc")
	}
	s := &Suite{
		DB:       db,
		Config:   config,
		BaseURL:  fmt.Sprintf("http://%s", config.ApiServerPort),
		GrpcConn: grpcConn,
	}

	return s
}
