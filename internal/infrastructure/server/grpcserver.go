package server

import (
	"fmt"
	"net"

	teamv1 "github.com/Junaidmdv/goalcircle-protos/team/v1"
	userclient_pb "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teaminvite"
	teammember_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	team_handler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/team"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	postgres "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/persistence/postgres"
	teamsaga "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	userclientcon "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/userclient"
	playeruc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/player"
	staffuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/staff"
	team_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	inviteuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teaminvite"
	teammemberuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teammember"

	inviteHandler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/invite"
	playerHandler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/player"
	staffHandler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/staff"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/validater"

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
	if err := psqldb.Migration(); err != nil {
		return err
	}
	validater, err := validater.NewValidater()
	if err != nil {
		return err
	}

	objectStorage,err:=storage.ObjectStorageFactoryMethod(*gs.Config.StorageConfig,gs.logger) 
	if err != nil{
		return err
	}
	codeGenerater := code.NewCodeGenerater(entity.CodeLength)



	teamRepository := team_repo.NewTeamRepository(psqldb.DB, gs.logger)
	TeamMemberRepository := teammember_repo.NewTeamMemberRepository(psqldb.DB, gs.logger)
	staffRepository := staffrepo.NewStaffRepository(psqldb.DB, gs.logger)
	playerRepo := player.NewPlayerRepository(psqldb.DB, gs.logger)
	teamInviteRepo := teaminvite.NewTeamMemberInviteRepository(psqldb.DB, gs.logger)

	teamUsecase := team_uc.NewTeamUsecase(teamRepository, gs.logger, codeGenerater, TeamMemberRepository, playerRepo,objectStorage,gs.Config.StorageConfig)
	teamMemberUc := teammemberuc.NewTeamMemberUsecase(TeamMemberRepository, teamInviteRepo, gs.logger)
	teamStaffUc := staffuc.NewTeamStaffUsecase(TeamMemberRepository, staffRepository)
	playerUc := playeruc.NewPlayerUsecase(playerRepo, TeamMemberRepository, teamRepository, codeGenerater)
	teaminviteUc := inviteuc.NewTeamInviteUsecase(teamInviteRepo)

	userclient, err := userclientcon.NewUserGRPCClient(gs.Config.UserSrv, gs.logger)
	if err != nil {
		return err
	}

	userAuthpb := userclient_pb.NewAuthServiceClient(userclient.Conn)

	teamSaga := teamsaga.NewTeamSagaMaker(teamUsecase, teamMemberUc, userAuthpb, gs.logger)
	teamHandler := team_handler.NewTeamHandler(teamUsecase, teamSaga, gs.Config.Server.TimeOut,validater)
	playerHandler := playerHandler.NewPlayerHandler(playerUc, gs.logger, &gs.Config.Server.TimeOut, validater)
	staffHandler := staffHandler.NewStaffHandler(teamStaffUc, gs.Config.Server.TimeOut)
	inviteHandler.NewTeamInviteHandler(teaminviteUc, gs.Config.Server.TimeOut)

	teamv1.RegisterTeamServiceServer(gs.Server, teamHandler)
	teamv1.RegisterPlayerServiceServer(gs.Server, playerHandler)
	teamv1.RegisterStaffServiceServer(gs.Server, staffHandler)

	return nil
}

func (gs *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gs.Config.Server.Port))

	if err != nil {

		return err
	}
	return gs.Server.Serve(lis)
}

func (gs *GRPCServer) GracefulShutdown() {
	gs.Server.GracefulStop()
}
