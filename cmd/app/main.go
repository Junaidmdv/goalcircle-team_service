package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/server"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

func main() {

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialise logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "logger sync error: %v\n", err)
		}
	}()

	cnfg, errs := config.LoadConfig().
		WithGrpcServer().
		WithPostgres().
		WithMinio().
		WithUserService().
		Build()

	if errs != nil {
		for _, err := range errs {
			logger.Error("configration error", "error", err)
			fmt.Println()
		}
		return
	}

	server := server.NewGRPCServer(cnfg, logger)
	if err := server.BootstrapSetup(); err != nil {
		logger.Error("bootstrap setup error", "error", err)
		return
	}

	go func() {
		logger.Info("grpc server running", "port", cnfg.Server.Port)
		if err := server.Run(); err != nil {
			logger.Error("grpc server error", "error", err)
			return
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	signl := <-signalChan
	logger.Info("gracefull shutdown signal recieved", "signal", signl.String())
	signal.Stop(signalChan)
	server.GracefulShutdown()

}
