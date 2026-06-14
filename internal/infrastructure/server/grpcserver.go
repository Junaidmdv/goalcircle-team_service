package server

import (
	"fmt"
	"net"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Server *grpc.Server
	Config *config.Config
	logger logger.Logger
}

func NewGRPCServer(cnfg *config.Config, logger logger.Logger) *GRPCServer {

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			RecoveryInterceptor(logger),
		),
	)

	return &GRPCServer{
		Server: server,
		Config: cnfg,
		logger: logger,
	}
}

func (gs *GRPCServer) BootstrapSetup() {

}

func (gs *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%d", gs.Config.Server.Port))
	if err != nil {
		return err
	}
	return gs.Server.Serve(lis)
}

func (gs *GRPCServer) GracefulShutdown() {
	gs.Server.GracefulStop()
}
