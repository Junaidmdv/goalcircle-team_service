package server

import (
	"fmt"
	"net"

	teampb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	userclient_pb "github.com/Junaidmdv/goalcircle-protos/user/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/player"
	playerinviterepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/playerinvite"
	staffrepo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staff"
	staff_invite_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/staffinvite"
	team_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/team"
	teammember_repo "github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teammember"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/teamplayer"
	"github.com/Junaidmdv/goalcircle-team_service/internal/domain/repository/transaction"
	playerinvitehandler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/invite"
	team_handler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/team"
	code "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/invitation"
	postgres "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/persistence/postgres"
	teamsaga "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/saga/team"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage"
	userclientcon "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/userclient"
	playerinviteuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/invite"
	playeruc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/player"
	staffuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/staff"
	team_uc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/team"
	teammemberuc "github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teammember"

	playerHandler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/player"
	staffHandler "github.com/Junaidmdv/goalcircle-team_service/internal/handler/grpc/staff"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/datetime"
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

	objectStorage, err := storage.ObjectStorageFactoryMethod(*gs.Config.StorageConfig, gs.logger)
	if err != nil {
		return err
	}
	codeGenerater := code.NewCodeGenerater(entity.CodeLength)

	datecalculater := datetime.NewDateCalculator()

	transactionHandler := transaction.NewUnitOfWork(psqldb.DB, gs.logger)

	teamRepository := team_repo.NewTeamRepository(psqldb.DB, gs.logger)
	TeamMemberRepository := teammember_repo.NewTeamMemberRepository(psqldb.DB, gs.logger)
	staffRepository := staffrepo.NewStaffRepository(psqldb.DB, gs.logger)
	playerRepo := player.NewPlayerRepository(psqldb.DB, gs.logger)
	staffInviteRepo := staff_invite_repo.NewStaffInviteRepository(*psqldb.DB, gs.logger)
	teamPlayerRepo := teamplayer.NewTeamPlayerRepository(psqldb.DB, gs.logger)
	playerInviteRepo := playerinviterepo.NewPlayerInviteRepository(psqldb.DB, gs.logger)

	teamUsecase := team_uc.NewTeamUsecase(teamRepository, gs.logger, codeGenerater, TeamMemberRepository, teamPlayerRepo, playerRepo, objectStorage, gs.Config.StorageConfig, datecalculater, transactionHandler, staffRepository)
	teamMemberUc := teammemberuc.NewTeamMemberUsecase(TeamMemberRepository, staffRepository, gs.logger, datecalculater)
	teamStaffUc := staffuc.NewTeamStaffUsecase(TeamMemberRepository, staffRepository, teamRepository, staffInviteRepo, teamPlayerRepo, datecalculater, objectStorage, gs.Config.StorageConfig, codeGenerater, gs.logger)
	playerUc := playeruc.NewPlayerUsecase(playerRepo, TeamMemberRepository, teamRepository, staffRepository, teamPlayerRepo, gs.logger, objectStorage, gs.Config.StorageConfig, codeGenerater)
	playerInviteUc := playerinviteuc.NewPlayerInviteUsecase(staffRepository, teamPlayerRepo, playerInviteRepo, playerRepo, objectStorage, TeamMemberRepository)
	userclient, err := userclientcon.NewUserGRPCClient(gs.Config.UserSrv, gs.logger)
	if err != nil {
		return err
	}

	userAuthpb := userclient_pb.NewAuthServiceClient(userclient.Conn)

	teamSaga := teamsaga.NewTeamSagaMaker(teamUsecase, teamMemberUc, userAuthpb, teamStaffUc, gs.logger)
	joinStaffsaga := teamsaga.NewJoinStaffSaga(teamStaffUc, userAuthpb, gs.logger)
	playersaga := teamsaga.NewPlayerSaga(playerUc, userAuthpb, gs.logger)

	teamHandler := team_handler.NewTeamHandler(teamUsecase, teamSaga, gs.Config.Server.TimeOut, validater)
	playerHandler := playerHandler.NewPlayerHandler(playerUc, gs.logger, &gs.Config.Server.TimeOut, validater, playersaga)
	staffHandler := staffHandler.NewStaffHandler(teamStaffUc, gs.Config.Server.TimeOut, validater, joinStaffsaga)
	playerInviteHandler := playerinvitehandler.NewPlayerInviteHandler(validater, playerInviteUc, gs.Config.Server.TimeOut)

	teampb.RegisterTeamServiceServer(gs.Server, teamHandler)
	teampb.RegisterPlayerServiceServer(gs.Server, playerHandler)
	teampb.RegisterStaffServiceServer(gs.Server, staffHandler)
	teampb.RegisterPlayerInviteServiceServer(gs.Server, playerInviteHandler)

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
