package invitehandler

import (
	"context"
	"time"

	pb "github.com/Junaidmdv/goalcircle-protos/team/v1"
	"github.com/Junaidmdv/goalcircle-team_service/internal/usecase/teaminvite"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TeamInviteHandler struct {
	teamInviteUsecase teaminvite.TeamInviteUsecase
	timeOut           time.Duration
	pb.UnimplementedTeamInviteServer
}

func NewTeamInviteHandler(ti teaminvite.TeamInviteUsecase, timeout time.Duration) *TeamInviteHandler {
	return &TeamInviteHandler{
		teamInviteUsecase: ti,
		timeOut:           timeout,
	}
}

func (th *TeamInviteHandler) CreateInvitation(ctx context.Context, req *pb.CreateInvitationReq) (*pb.CreateInvitationRes, error) {
	context, cancel := context.WithTimeout(ctx, th.timeOut)
	defer cancel()

	res, err := th.teamInviteUsecase.CreateInvitation(context, &teaminvite.TeamInviteReq{
		TeamMemberID: req.TeamMemberId,
		Role:         req.Role,
	})

	if err != nil {
		return nil, apperror.GRPCStatus(err)
	}

	return &pb.CreateInvitationRes{
		Code:      res.Code,
		ExpiredAt: timestamppb.New(res.ExpiresAt),
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
}
