package server

import (
	"fmt"
	"net"

	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	userclient_pb "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	teammember_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	team_handler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/team"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	postgres "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/persistence/postgres"
	teamsaga "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga/team"
	userclientcon "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/userclient"
	team_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	teammember_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teamowner"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
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

func (gs *GRPCServer) BootstrapSetup() error {

	psqldb, err := postgres.NewPostgresDB(gs.Config.Postgres)
	if err != nil {
		return err
	} 

	codeGenerater:=code.NewCodeGenerater(entity.CodeLength)
	teamRepository := team_repo.NewTeamRepository(psqldb.DB, gs.logger)
	TeamMemberRepository := teammember_repo.NewTeamMemberRepository(psqldb.DB, gs.logger)

	teamUsecase := team_uc.NewTeamUsecase(teamRepository, gs.logger,codeGenerater)
	teamMemberUc := teammember_uc.NewTeamOwnerUsecase(TeamMemberRepository)

	userclient, err := userclientcon.NewUserGRPCClient(gs.Config.UserSrv, gs.logger)
	if err != nil {
		return err
	}

	userAuthpb := userclient_pb.NewAuthServiceClient(userclient.Conn)

	teamSaga := teamsaga.NewTeamSagaMaker(teamUsecase, teamMemberUc, userAuthpb, gs.logger)
	teamHandler := team_handler.NewTeamHandler(teamUsecase, teamSaga, gs.Config.Server.TimeOut)
	teamv1.RegisterTeamServiceServer(gs.Server, teamHandler)

	return nil
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
