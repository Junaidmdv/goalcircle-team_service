package server

import (
	"fmt"
	"net"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"

	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	team_handler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/team"
	team_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
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
	teamUsecase := team_uc.NewTeamUsecase()
	teamHandler := team_handler.NewTeamHandler(teamUsecase)  
	





	teamv1.RegisterTeamServiceServer(gs.Server, teamHandler)

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
