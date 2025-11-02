package bookgrpc

import (
	"context"
	analyticsv1 "github.com/s-khechnev/pet-project/protos/gen/go/analytics"
	"google.golang.org/grpc"
)

type serverApi struct {
	analyticsv1.UnimplementedAnalyticsServer
}

func Register(server *grpc.Server) {
	analyticsv1.RegisterAnalyticsServer(server, &serverApi{})
}

func (s *serverApi) GetStatistics(
	ctx context.Context,
	request *analyticsv1.StatisticsRequest,
) (*analyticsv1.StatisticsResponse, error) {
	return &analyticsv1.StatisticsResponse{
		CountTextSymbols: 123333333333333,
	}, nil
}
