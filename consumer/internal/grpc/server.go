package bookgrpc

import (
	"consumer/internal/service/analytics"
	"context"
	analyticsv1 "github.com/s-khechnev/pet-project/protos/gen/go/analytics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerApi struct {
	analyticsService *analytics.BookAnalyticsService
	analyticsv1.UnimplementedAnalyticsServer
}

func NewServerApi(analyticsService *analytics.BookAnalyticsService) *ServerApi {
	return &ServerApi{
		analyticsService: analyticsService,
	}
}

func Register(server *grpc.Server, api *ServerApi) {
	analyticsv1.RegisterAnalyticsServer(server, api)
}

func (s *ServerApi) GetStatistics(
	ctx context.Context,
	_ *analyticsv1.StatisticsRequest,
) (*analyticsv1.StatisticsResponse, error) {

	stats, err := s.analyticsService.GetStatistics(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal erorr: %v", err.Error())
	}

	return &analyticsv1.StatisticsResponse{
		CountTextSymbols: stats.CountTextSymbols,
		CountBooks:       stats.CountBooks,
		CountAuthors:     stats.CountAuthors,
	}, nil
}
