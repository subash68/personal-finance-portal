package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/subash68/ate/ate_token_service/configuration"
	"github.com/subash68/ate/ate_token_service/pkg/logger"
	"github.com/subash68/ate/ate_token_service/pkg/protocol/grpc"
	"github.com/subash68/ate/ate_token_service/pkg/protocol/rest"
	v1 "github.com/subash68/ate/ate_token_service/pkg/service/v1"
)

type Config struct {
	GRPCPort string
	HTTPPort string

	LogLevel      int
	LogTimeFormat string

	// DatabaseInstance string
	// DatabaseName string
	// Collection string
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config

	configuration.LoadConfig()

	cfg.GRPCPort = strconv.Itoa(configuration.PortConfig().GRPCPort)
	cfg.HTTPPort = strconv.Itoa(configuration.PortConfig().HTTPPort)

	cfg.LogLevel = configuration.LogConfig().LogLevel
	cfg.LogTimeFormat = configuration.LogConfig().LogTimeFormat

	// cfg.DatabaseInstance = configuration.DbConfig().DatabaseInstance
	// cfg.DatabaseName = configuration.DbConfig().DatabaseName
	// cfg.Collection = configuration.DbConfig().Collection

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid tcp port for grpc server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid http port for http gateway: '%s'", cfg.HTTPPort)
	}

	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	//GET mongo db isntance for parameter

	// db, err := database.GetMongoClient()
	// if err != nil {
	// 	return fmt.Errorf("error while connecting with database %v", err)
	// }

	v1API := v1.NewTokenServiceServer()

	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
