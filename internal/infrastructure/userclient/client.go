package userclient

import (
	"fmt"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCUserClient struct {
	Conn *grpc.ClientConn
}

func NewUserGRPCClient(serviceCnfg *config.UserService, logger logger.Logger) (*GRPCUserClient, error) {
	addr := fmt.Sprintf("%s:%d", serviceCnfg.Host, serviceCnfg.Port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %d: %w", serviceCnfg.Port, err)
	}
	logger.Info("grpc connected", "host", serviceCnfg.Host, "port", serviceCnfg.Port)
	return &GRPCUserClient{Conn: conn}, nil
}

func (c *GRPCUserClient) Close() error {
	return c.Conn.Close()
}

